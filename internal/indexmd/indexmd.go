// Package indexmd renders the instance index.md view.
package indexmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertguss/mycelium/internal/idpath"
	"github.com/robertguss/mycelium/internal/logfmt"
	"github.com/robertguss/mycelium/internal/manifest"
)

const maxLogTail = 20

// NSCount is one registered namespace's artifact count.
type NSCount struct {
	NS    string
	Count int
}

// Instance is the data needed to render index.md.
type Instance struct {
	IdeaName   string
	State      string
	Tier       string
	Revisit    string
	GithubRepo string
	Counts     []NSCount
	LogLines   []string // all parseable lines; Render keeps the last 20
	Wake       string   // "none" or "briefs/LATEST.md"
}

// ZeroCounts returns a Counts slice for every registered NS at zero.
func ZeroCounts() []NSCount {
	types := idpath.Types()
	out := make([]NSCount, len(types))
	for i, t := range types {
		out[i] = NSCount{NS: t.NS}
	}
	return out
}

// Inc increments the count for ns (no-op if unknown).
func (inst *Instance) Inc(ns string) {
	for i := range inst.Counts {
		if inst.Counts[i].NS == ns {
			inst.Counts[i].Count++
			return
		}
	}
}

// Render returns deterministic index.md bytes for inst.
func Render(inst Instance) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", inst.IdeaName)

	fmt.Fprintf(&b, "## State\n\n")
	fmt.Fprintf(&b, "state: %s\n", inst.State)
	fmt.Fprintf(&b, "tier: %s\n", inst.Tier)
	fmt.Fprintf(&b, "revisit: %s\n", inst.Revisit)
	fmt.Fprintf(&b, "github_repo: %s\n\n", inst.GithubRepo)

	fmt.Fprintf(&b, "## Artifacts\n\n")
	counts := inst.Counts
	if len(counts) == 0 {
		counts = ZeroCounts()
	}
	for _, c := range counts {
		fmt.Fprintf(&b, "%s: %d\n", c.NS, c.Count)
	}
	b.WriteByte('\n')

	fmt.Fprintf(&b, "## Log tail\n\n")
	tail := inst.LogLines
	if len(tail) > maxLogTail {
		tail = tail[len(tail)-maxLogTail:]
	}
	for _, line := range tail {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	b.WriteByte('\n')

	wake := inst.Wake
	if wake == "" {
		wake = "none"
	}
	fmt.Fprintf(&b, "## Wake\n\n%s\n", wake)
	return []byte(b.String())
}

// Load reads an instance root into an Instance.
func Load(root string) (Instance, error) {
	mb, err := os.ReadFile(filepath.Join(root, "mycelium.toml"))
	if err != nil {
		return Instance{}, err
	}
	m, err := manifest.Parse(mb)
	if err != nil {
		return Instance{}, err
	}
	logBytes, err := os.ReadFile(filepath.Join(root, "log.md"))
	if err != nil {
		return Instance{}, err
	}
	counts, err := countArtifacts(root)
	if err != nil {
		return Instance{}, err
	}
	wake := "none"
	if _, err := os.Stat(filepath.Join(root, "briefs", "LATEST.md")); err == nil {
		wake = "briefs/LATEST.md"
	}
	return Instance{
		IdeaName:   m.IdeaName,
		State:      m.State,
		Tier:       m.Tier,
		Revisit:    m.Revisit,
		GithubRepo: m.GithubRepo,
		Counts:     counts,
		LogLines:   logfmt.ParseableLines(logBytes),
		Wake:       wake,
	}, nil
}

func countArtifacts(root string) ([]NSCount, error) {
	out := ZeroCounts()
	byNS := map[string]int{}
	for i, c := range out {
		byNS[c.NS] = i
	}
	for _, t := range idpath.Types() {
		dir := filepath.Join(root, t.Home)
		st, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		if !st.IsDir() {
			continue
		}
		err = filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				return nil
			}
			if d.Name() == "README.md" {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			id, _, err := idpath.ParsePath(rel)
			if err != nil {
				return nil
			}
			if i, ok := byNS[id.NS]; ok {
				out[i].Count++
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
