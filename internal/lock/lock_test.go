package lock_test

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/robertguss/mycelium/internal/lock"
)

func TestAcquireRelease(t *testing.T) {
	root := t.TempDir()
	h, err := lock.Acquire(root, time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(lock.Path(root))
	if err != nil {
		t.Fatal(err)
	}
	want := "pid="
	if len(b) < 4 || string(b[:4]) != want {
		t.Fatalf("lock content = %q", b)
	}
	if err := h.Release(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(lock.Path(root)); !os.IsNotExist(err) {
		t.Fatalf("lock file should be removed, err=%v", err)
	}
}

func TestLiveLockRefuses(t *testing.T) {
	root := t.TempDir()
	h, err := lock.Acquire(root, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	defer h.Release()

	_, err = lock.Acquire(root, time.Now())
	if !errors.Is(err, lock.ErrBusy) {
		t.Fatalf("want ErrBusy, got %v", err)
	}
}

func TestStaleDeadPID(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".mycelium")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	dead := "pid=2147483646\nstarted=2020-01-01T00:00:00Z\n"
	if err := os.WriteFile(lock.Path(root), []byte(dead), 0o644); err != nil {
		t.Fatal(err)
	}
	info, err := lock.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if info.State != lock.Stale {
		t.Fatalf("state = %v, want Stale", info.State)
	}
	if err := lock.RemoveStale(root); err != nil {
		t.Fatal(err)
	}
	info, err = lock.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if info.State != lock.Absent {
		t.Fatalf("after remove: %v", info.State)
	}
}

func TestStaleMissingPID(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".mycelium"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lock.Path(root), []byte("started=2020-01-01T00:00:00Z\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	info, err := lock.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if info.State != lock.Stale {
		t.Fatalf("state = %v, want Stale", info.State)
	}
}

func TestInspectLive(t *testing.T) {
	root := t.TempDir()
	h, err := lock.Acquire(root, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	defer h.Release()
	info, err := lock.Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if info.State != lock.Live || info.PID != h.PID() {
		t.Fatalf("info = %+v", info)
	}
}

func TestAcquireReusesLeftoverInode(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".mycelium"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := lock.Path(root)
	if err := os.WriteFile(p, []byte("pid=2147483646\nstarted=2020-01-01T00:00:00Z\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	inoBefore := fileIno(t, p)
	h, err := lock.Acquire(root, time.Date(2026, 8, 15, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	defer h.Release()
	inoAfter := fileIno(t, p)
	if inoBefore != inoAfter {
		t.Fatalf("inode changed: before=%d after=%d (Acquire must flock existing file)", inoBefore, inoAfter)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), fmt.Sprintf("pid=%d", h.PID())) {
		t.Fatalf("pid not overwritten: %q", b)
	}
}

// H1 + M3: two processes race Acquire against a leftover dead-pid lock file;
// exactly one Held, the other ErrBusy.
func TestTwoProcOneHeldAgainstDeadPID(t *testing.T) {
	if os.Getenv("MYCELIUM_LOCK_CHILD") == "1" {
		root := os.Getenv("MYCELIUM_LOCK_ROOT")
		h, err := lock.Acquire(root, time.Now())
		if err != nil {
			fmt.Println("BUSY")
			os.Exit(2)
		}
		fmt.Println("HELD")
		// Hold until parent closes stdin.
		_, _ = bufio.NewReader(os.Stdin).ReadByte()
		_ = h.Release()
		os.Exit(0)
	}

	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".mycelium"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lock.Path(root), []byte("pid=2147483646\nstarted=2020-01-01T00:00:00Z\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	inoBefore := fileIno(t, lock.Path(root))

	startChild := func() (*exec.Cmd, *bufio.Reader, io.WriteCloser) {
		t.Helper()
		cmd := exec.Command(os.Args[0], "-test.run=^TestTwoProcOneHeldAgainstDeadPID$", "-test.v=false")
		cmd.Env = append(os.Environ(),
			"MYCELIUM_LOCK_CHILD=1",
			"MYCELIUM_LOCK_ROOT="+root,
		)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatal(err)
		}
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		return cmd, bufio.NewReader(stdout), stdin
	}

	// First child should win the flock on the leftover inode.
	cmd1, r1, in1 := startChild()
	line1, err := r1.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(line1) != "HELD" {
		t.Fatalf("child1 = %q, want HELD", line1)
	}

	cmd2, r2, in2 := startChild()
	line2, err := r2.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(line2) != "BUSY" {
		t.Fatalf("child2 = %q, want BUSY (exactly one Held)", line2)
	}
	_ = in2.Close()
	_ = cmd2.Wait()

	inoAfter := fileIno(t, lock.Path(root))
	if inoBefore != inoAfter {
		t.Fatalf("inode split: before=%d after=%d", inoBefore, inoAfter)
	}

	_ = in1.Close()
	if err := cmd1.Wait(); err != nil {
		t.Fatal(err)
	}
}

func TestTwoProcLiveLockBusy(t *testing.T) {
	if os.Getenv("MYCELIUM_LOCK_CHILD") == "hold" {
		root := os.Getenv("MYCELIUM_LOCK_ROOT")
		h, err := lock.Acquire(root, time.Now())
		if err != nil {
			fmt.Println("ERR", err)
			os.Exit(1)
		}
		fmt.Println("HELD")
		_, _ = bufio.NewReader(os.Stdin).ReadByte()
		_ = h.Release()
		os.Exit(0)
	}
	if os.Getenv("MYCELIUM_LOCK_CHILD") == "probe" {
		root := os.Getenv("MYCELIUM_LOCK_ROOT")
		_, err := lock.Acquire(root, time.Now())
		if errors.Is(err, lock.ErrBusy) {
			fmt.Println("BUSY")
			os.Exit(0)
		}
		fmt.Println("UNEXPECTED", err)
		os.Exit(1)
	}

	root := t.TempDir()
	holder := exec.Command(os.Args[0], "-test.run=^TestTwoProcLiveLockBusy$", "-test.v=false")
	holder.Env = append(os.Environ(), "MYCELIUM_LOCK_CHILD=hold", "MYCELIUM_LOCK_ROOT="+root)
	stdout, err := holder.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdin, err := holder.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	holder.Stderr = os.Stderr
	if err := holder.Start(); err != nil {
		t.Fatal(err)
	}
	line, err := bufio.NewReader(stdout).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(line) != "HELD" {
		t.Fatalf("holder = %q", line)
	}

	probe := exec.Command(os.Args[0], "-test.run=^TestTwoProcLiveLockBusy$", "-test.v=false")
	probe.Env = append(os.Environ(), "MYCELIUM_LOCK_CHILD=probe", "MYCELIUM_LOCK_ROOT="+root)
	out, err := probe.CombinedOutput()
	if err != nil {
		t.Fatalf("probe: %v %s", err, out)
	}
	if !strings.Contains(string(out), "BUSY") {
		t.Fatalf("probe out = %q", out)
	}

	_ = stdin.Close()
	_ = holder.Wait()
}

func fileIno(t *testing.T, path string) uint64 {
	t.Helper()
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		t.Fatal("stat Sys is not *syscall.Stat_t")
	}
	return st.Ino
}
