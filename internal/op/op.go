// Package op implements the lock/journal/stage/commit/rollback protocol.
package op

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/robertguss/mycelium/internal/journal"
	"github.com/robertguss/mycelium/internal/lock"
)

var (
	ErrJournalMismatch = errors.New("op: leftover journal for a different command")
	ErrNothingToAbort  = errors.New("op: nothing to abort")
	ErrPartialCommit   = errors.New("op: cannot rollback after partial commit")
	ErrLocked          = errors.New("op: lock held")
)

// Intent describes a mutating operation before staging.
type Intent struct {
	Op         string
	Type       string // empty → JSON null
	Title      string
	OriginalID string
	LogLine    string
	Argv       []string
	OpID       string // stage subdirectory name; empty → derived
}

// Staged is one file written under the stage dir, destined for RelTo.
type Staged struct {
	RelTo   string // destination path relative to root
	Content []byte
}

// Session is one in-flight protocol run (holding the lock).
type Session struct {
	root    string
	held    *lock.Held
	journal *journal.Journal
}

// Begin preflights leftover journal, acquires the lock, and prepares a session.
// If a matching leftover journal exists, OriginalID and undone renames are reused.
func Begin(root string, intent Intent, now time.Time) (*Session, error) {
	existing, err := journal.Load(root)
	if err != nil && !errors.Is(err, journal.ErrNotExist) {
		return nil, err
	}
	if existing != nil {
		if !journal.Matches(existing, intent.Argv, intent.Op, intent.Type, intent.Title) {
			return nil, fmt.Errorf("%w: use mycelium check --abort-journal", ErrJournalMismatch)
		}
		intent.OriginalID = existing.OriginalID
	}

	// Acquire opens the existing lock path (O_CREATE|O_RDWR) and flocks that
	// inode. Do not RemoveStale first: unlinking then recreating splits inodes.
	held, err := lock.Acquire(root, now)
	if err != nil {
		if errors.Is(err, lock.ErrBusy) {
			return nil, fmt.Errorf("%w: %v", ErrLocked, err)
		}
		return nil, err
	}

	opID := intent.OpID
	if opID == "" {
		opID = fmt.Sprintf("%s-%d", intent.Op, now.UTC().UnixNano())
	}

	s := &Session{root: root, held: held}
	if existing != nil {
		s.journal = existing
		return s, nil
	}

	j := &journal.Journal{
		SchemaVersion: 1,
		Op:            intent.Op,
		Title:         intent.Title,
		OriginalID:    intent.OriginalID,
		StartedAt:     now.UTC().Format(time.RFC3339),
		StagedDir:     filepath.ToSlash(filepath.Join(".mycelium", "stage", opID)),
		LogLine:       intent.LogLine,
		Argv:          append([]string(nil), intent.Argv...),
	}
	j.SetType(intent.Type)
	s.journal = j
	return s, nil
}

// OriginalID returns the ID recorded for this session (resume-stable).
func (s *Session) OriginalID() string {
	if s == nil || s.journal == nil {
		return ""
	}
	return s.journal.OriginalID
}

// Journal returns the in-memory journal (for tests and callers).
func (s *Session) Journal() *journal.Journal {
	return s.journal
}

// StageDir returns the absolute stage directory.
func (s *Session) StageDir() string {
	return filepath.Join(s.root, filepath.FromSlash(s.journal.StagedDir))
}

// Stage writes files into the stage dir and records rename intents.
// RelTo paths are relative to the instance root. Call once before Commit
// on a fresh session; on resume, Stage is a no-op if renames already exist.
func (s *Session) Stage(files []Staged) error {
	if s.journal == nil {
		return errors.New("op: no journal")
	}
	if len(s.journal.Renames) > 0 {
		// Resume: staged files and renames already recorded.
		return nil
	}
	stageAbs := s.StageDir()
	if err := os.MkdirAll(stageAbs, 0o755); err != nil {
		return err
	}
	renames := make([]journal.Rename, 0, len(files))
	for i, f := range files {
		if f.RelTo == "" {
			return fmt.Errorf("op: empty RelTo at %d", i)
		}
		base := filepath.Base(f.RelTo)
		if base == "." || base == "/" || base == "" {
			base = fmt.Sprintf("file-%d", i)
		}
		stagedName := fmt.Sprintf("%03d-%s", i, base)
		stagedRel := filepath.ToSlash(filepath.Join(s.journal.StagedDir, stagedName))
		abs := filepath.Join(s.root, filepath.FromSlash(stagedRel))
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(abs, f.Content, 0o644); err != nil {
			return err
		}
		renames = append(renames, journal.Rename{
			From: stagedRel,
			To:   filepath.ToSlash(f.RelTo),
			Done: false,
		})
	}
	s.journal.Renames = renames
	return journal.Save(s.root, s.journal)
}

// SetOriginalID updates the journal's original_id before first commit.
func (s *Session) SetOriginalID(id string) error {
	s.journal.OriginalID = id
	return journal.Save(s.root, s.journal)
}

// Commit applies all undone renames in order, persisting the journal after each.
// On full success: delete journal, delete staged dir, release lock, remove lock file.
func (s *Session) Commit() error {
	_, err := s.commitN(-1)
	return err
}

// CommitPartial applies at most n undone renames and leaves the journal.
// Used by tests to simulate crash-after-partial-rename. n < 0 means all.
func (s *Session) CommitPartial(n int) error {
	_, err := s.commitN(n)
	return err
}

func (s *Session) commitN(max int) (int, error) {
	applied := 0
	for i := range s.journal.Renames {
		if s.journal.Renames[i].Done {
			continue
		}
		if max >= 0 && applied >= max {
			if err := journal.Save(s.root, s.journal); err != nil {
				return applied, err
			}
			return applied, nil
		}
		from := filepath.Join(s.root, filepath.FromSlash(s.journal.Renames[i].From))
		to := filepath.Join(s.root, filepath.FromSlash(s.journal.Renames[i].To))
		_, fromErr := os.Stat(from)
		_, toErr := os.Stat(to)
		fromOK := fromErr == nil
		toOK := toErr == nil
		if toOK && !fromOK {
			// Crash between Rename and Save: dest landed, source gone.
			s.journal.Renames[i].Done = true
			applied++
			if err := journal.Save(s.root, s.journal); err != nil {
				return applied, err
			}
			continue
		}
		if toOK && fromOK {
			_ = journal.Save(s.root, s.journal)
			return applied, fmt.Errorf("op: rename conflict: both %s and %s exist", from, to)
		}
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			_ = journal.Save(s.root, s.journal)
			return applied, err
		}
		if err := os.Rename(from, to); err != nil {
			_ = journal.Save(s.root, s.journal)
			return applied, err
		}
		s.journal.Renames[i].Done = true
		applied++
		if err := journal.Save(s.root, s.journal); err != nil {
			return applied, err
		}
	}
	allDone := true
	for _, r := range s.journal.Renames {
		if !r.Done {
			allDone = false
			break
		}
	}
	if allDone && max < 0 {
		return applied, s.finish()
	}
	return applied, nil
}

func (s *Session) finish() error {
	staged := s.StageDir()
	_ = os.RemoveAll(staged)
	// prune empty stage parent if possible
	_ = os.Remove(filepath.Join(s.root, ".mycelium", "stage"))
	if err := journal.Remove(s.root); err != nil {
		return err
	}
	return s.held.Release()
}

// Rollback deletes staged files and journal when no rename has completed.
// After a partial commit, returns ErrPartialCommit and leaves the journal.
func (s *Session) Rollback() error {
	for _, r := range s.journal.Renames {
		if r.Done {
			return ErrPartialCommit
		}
	}
	_ = os.RemoveAll(s.StageDir())
	_ = os.Remove(filepath.Join(s.root, ".mycelium", "stage"))
	if err := journal.Remove(s.root); err != nil {
		return err
	}
	return s.held.Release()
}

// Close releases the lock without finishing the journal. Use after CommitPartial
// to simulate a crash that left the journal on disk.
func (s *Session) Close() error {
	if s == nil || s.held == nil {
		return nil
	}
	err := s.held.Release()
	s.held = nil
	return err
}

// Abort clears staged temps, journal, and a stale lock. Does not delete
// already-renamed destinations. Returns ErrNothingToAbort when neither
// journal nor stale lock exists. A live lock refuses immediately.
func Abort(root string) error {
	info, err := lock.Inspect(root)
	if err != nil {
		return err
	}
	if info.State == lock.Live {
		return fmt.Errorf("%w: pid=%d", ErrLocked, info.PID)
	}

	j, err := journal.Load(root)
	hasJournal := err == nil
	if err != nil && !errors.Is(err, journal.ErrNotExist) {
		return err
	}
	stale := info.State == lock.Stale

	if !hasJournal && !stale {
		return ErrNothingToAbort
	}

	if hasJournal {
		for _, r := range j.Renames {
			if r.Done {
				continue
			}
			from := filepath.Join(root, filepath.FromSlash(r.From))
			_ = os.Remove(from)
		}
		if j.StagedDir != "" {
			_ = os.RemoveAll(filepath.Join(root, filepath.FromSlash(j.StagedDir)))
		}
		_ = os.Remove(filepath.Join(root, ".mycelium", "stage"))
		if err := journal.Remove(root); err != nil {
			return err
		}
	}
	if stale {
		if err := lock.RemoveStale(root); err != nil {
			return err
		}
	}
	return nil
}

// Detect reports leftover journal and/or stale lock for check.
func Detect(root string) (hasJournal bool, staleLock bool, err error) {
	_, jerr := journal.Load(root)
	if jerr == nil {
		hasJournal = true
	} else if !errors.Is(jerr, journal.ErrNotExist) {
		return false, false, jerr
	}
	info, err := lock.Inspect(root)
	if err != nil {
		return hasJournal, false, err
	}
	return hasJournal, info.State == lock.Stale, nil
}
