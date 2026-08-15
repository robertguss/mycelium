package lock_test

import (
	"errors"
	"os"
	"path/filepath"
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
	// PID 1<<22 is almost never a live process on Linux; use a high unused pid.
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
