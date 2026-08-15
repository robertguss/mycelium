// Package lock provides the exclusive .mycelium/lock flock.
package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const relPath = ".mycelium/lock"

var (
	ErrBusy  = errors.New("lock: held by another process")
	ErrStale = errors.New("lock: stale")
)

// State of the lock file on disk (not this process's hold).
type State int

const (
	Absent State = iota
	Live
	Stale
)

// Info describes an on-disk lock file.
type Info struct {
	State   State
	PID     int
	Started string
}

// Held is an acquired exclusive lock. Release unlocks and removes the file.
type Held struct {
	root string
	f    *os.File
	pid  int
}

// Path returns the absolute lock file path under root.
func Path(root string) string {
	return filepath.Join(root, relPath)
}

// Acquire takes a non-blocking exclusive flock. On busy, returns ErrBusy
// wrapping a message that includes the holder pid when known.
func Acquire(root string, now time.Time) (*Held, error) {
	dir := filepath.Join(root, ".mycelium")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	p := Path(root)
	f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		info, _ := Inspect(root)
		_ = f.Close()
		if info.State == Live {
			return nil, fmt.Errorf("%w: pid=%d", ErrBusy, info.PID)
		}
		return nil, fmt.Errorf("%w: flock busy", ErrBusy)
	}
	pid := os.Getpid()
	started := now.UTC().Format(time.RFC3339)
	content := fmt.Sprintf("pid=%d\nstarted=%s\n", pid, started)
	if err := f.Truncate(0); err != nil {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
		return nil, err
	}
	if _, err := f.Seek(0, 0); err != nil {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
		return nil, err
	}
	if _, err := f.WriteString(content); err != nil {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
		return nil, err
	}
	if err := f.Sync(); err != nil {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
		return nil, err
	}
	return &Held{root: root, f: f, pid: pid}, nil
}

// Release unlocks and removes the lock file.
func (h *Held) Release() error {
	if h == nil || h.f == nil {
		return nil
	}
	_ = syscall.Flock(int(h.f.Fd()), syscall.LOCK_UN)
	err := h.f.Close()
	h.f = nil
	_ = os.Remove(Path(h.root))
	return err
}

// PID returns the pid written into the lock file by this holder.
func (h *Held) PID() int {
	if h == nil {
		return 0
	}
	return h.pid
}

// Inspect reports whether a lock file is absent, live, or stale.
func Inspect(root string) (Info, error) {
	p := Path(root)
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return Info{State: Absent}, nil
		}
		return Info{}, err
	}
	pid, started := parseLock(string(b))
	if pid <= 0 {
		return Info{State: Stale, Started: started}, nil
	}
	if alive(pid) {
		return Info{State: Live, PID: pid, Started: started}, nil
	}
	return Info{State: Stale, PID: pid, Started: started}, nil
}

// RemoveStale deletes a stale lock file. Refuses if the lock is live.
func RemoveStale(root string) error {
	info, err := Inspect(root)
	if err != nil {
		return err
	}
	switch info.State {
	case Absent:
		return nil
	case Live:
		return fmt.Errorf("%w: pid=%d", ErrBusy, info.PID)
	case Stale:
		return os.Remove(Path(root))
	default:
		return fmt.Errorf("lock: unknown state %d", info.State)
	}
}

func parseLock(s string) (pid int, started string) {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "pid=") {
			pid, _ = strconv.Atoi(strings.TrimPrefix(line, "pid="))
		}
		if strings.HasPrefix(line, "started=") {
			started = strings.TrimPrefix(line, "started=")
		}
	}
	return pid, started
}

func alive(pid int) bool {
	// syscall.Kill(pid, 0) is the durable form of FindProcess+Signal(0).
	// os.Process.Signal maps never-started pids to ErrProcessDone, not ESRCH.
	err := syscall.Kill(pid, 0)
	if err == nil {
		return true
	}
	if errors.Is(err, syscall.ESRCH) {
		return false
	}
	return true
}
