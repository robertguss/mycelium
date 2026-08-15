package clitest_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/robertguss/mycelium/internal/clitest"
)

const ms401AppendixBCMP = `+++
id = "CMP-001"
title = "SQLite second opinion"
date = "2026-08-15"
rung = "second-opinion"
opt_in = true
cost_class = "cheap"
adapter = "manual"
+++

# CMP-001 — SQLite second opinion

<!-- slug: sqlite-second-opinion -->

## Prompt

Should this idea use SQLite as the store? Answer independently. Do not see other reports.

## Attachments

none

## Cost

cheap
`

const ms401AppendixBRPT = `+++
id = "RPT-001"
title = "Model A"
date = "2026-08-15"
model = "model-a"
commissioning = "CMP-001"
rung = "second-opinion"
adapter = "manual"
prompt_sha256 = "ec87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de"
+++

# RPT-001 — Model A

<!-- slug: model-a -->

## Position

Use SQLite.

## Findings

One file. Enough.

## Dissent

none
`

const ms401SecondOpinionRPT2 = `+++
id = "RPT-002"
title = "Model B"
date = "2026-08-15"
model = "model-b"
commissioning = "CMP-001"
rung = "second-opinion"
adapter = "manual"
prompt_sha256 = "ec87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de"
+++

# RPT-002 — Model B

<!-- slug: model-b -->

## Position

Do not use SQLite.

## Findings

A second writer changes the tradeoff.

## Dissent

none
`

const ms401AppendixCCMP = `+++
id = "CMP-001"
title = "SQLite council"
date = "2026-08-15"
rung = "council"
opt_in = true
cost_class = "standard"
adapter = "cursor"
+++

# CMP-001 — SQLite council

<!-- slug: sqlite-council -->

## Prompt

Review the SQLite store decision. Work independently. Do not see other reports. Retain dissent.

## Attachments

none

## Cost

standard
`

const ms401AppendixCRPT1 = `+++
id = "RPT-001"
title = "Model A"
date = "2026-08-15"
model = "model-a"
commissioning = "CMP-001"
rung = "council"
adapter = "cursor"
prompt_sha256 = "8997334f7f2f0bf821bce8ccc4a8d6cf027317c6c66d821200a032a6a11ce098"
+++

# RPT-001 — Model A

<!-- slug: model-a -->

## Position

Use SQLite.

## Findings

One file. Enough.

## Dissent

SEED-DISSENT
`

const ms401AppendixCRPT2 = `+++
id = "RPT-002"
title = "Model B"
date = "2026-08-15"
model = "model-b"
commissioning = "CMP-001"
rung = "council"
adapter = "cursor"
prompt_sha256 = "8997334f7f2f0bf821bce8ccc4a8d6cf027317c6c66d821200a032a6a11ce098"
+++

# RPT-002 — Model B

<!-- slug: model-b -->

## Position

Do not use SQLite.

## Findings

A second writer changes the tradeoff.

## Dissent

none
`

const ms401AppendixCRCL = `+++
id = "RCL-001"
title = "SQLite council reconciliation"
date = "2026-08-15"
commissioning = "CMP-001"
rung = "council"
+++

# RCL-001 — SQLite council reconciliation

<!-- slug: sqlite-council-reconciliation -->

## Convergence

Both reports address the store.

## Material disagreement

Whether SQLite is enough.

## Evidence unique to one report

none

## Contradictory evidence

none

## Different assumptions

Single-process vs later writers.

## Different scope interpretations

none

## Recommendations independently supported

Spike if a second writer appears.

## Questions requiring another spike

none

## Final reconciled recommendation

Use SQLite now. Revisit on a second writer.

## Retained dissent

SEED-DISSENT
`

type ms401Fixture struct {
	bin      string
	work     string
	env      []string
	home     string
	ghMarker string
}

func TestMS401HermeticFixtures(t *testing.T) {
	bin := clitest.Bin(t)
	stubDir, ghMarker := installGhNeverStub(t)
	home := t.TempDir()
	work := t.TempDir()
	env := append(
		hermeticEnv(stubDir, home),
		"GH_TOKEN=",
		"MYCELIUM_CONFIG="+filepath.Join(home, "config"),
	)
	f := ms401Fixture{
		bin:      bin,
		work:     work,
		env:      env,
		home:     home,
		ghMarker: ghMarker,
	}

	t.Run("enable-pack-present-no-reviews", func(t *testing.T) {
		inst := f.scaffold(t, "pack-on", "Pack On")
		if info, err := os.Stat(filepath.Join(inst, "program", "packs", "council")); err != nil || !info.IsDir() {
			t.Fatalf("council pack missing: %v", err)
		}
		ms401AssertAbsent(t, filepath.Join(inst, "reviews"))
		f.check(t, inst, 0)
	})

	t.Run("disable-pack-no-reviews", func(t *testing.T) {
		inst := f.scaffold(t, "pack-off", "Pack Off")
		ms401RemoveCouncilPack(t, inst)
		ms401AssertAbsent(t, filepath.Join(inst, "reviews"))
		f.check(t, inst, 0)
	})

	t.Run("disable-pack-with-reviews", func(t *testing.T) {
		inst := f.scaffold(t, "pack-off-reviews", "Pack Off Reviews")
		ms401RemoveCouncilPack(t, inst)
		if err := os.Mkdir(filepath.Join(inst, "reviews"), 0o755); err != nil {
			t.Fatal(err)
		}
		stderr := f.check(t, inst, 1)
		if !strings.Contains(stderr, "extra-top-level") || !strings.Contains(stderr, "reviews/") {
			t.Fatalf("stderr must identify reviews/ as extra-top-level: %q", stderr)
		}
	})

	t.Run("disable-pack-with-reviews-deviation", func(t *testing.T) {
		inst := f.scaffold(t, "pack-off-reviews-deviation", "Pack Off Reviews Deviation")
		ms401RemoveCouncilPack(t, inst)
		if err := os.Mkdir(filepath.Join(inst, "reviews"), 0o755); err != nil {
			t.Fatal(err)
		}
		addDeviation(t, inst, "extra-top-level:reviews/", "retain reviews after disabling the council pack")
		f.check(t, inst, 0)
	})

	t.Run("pack-namespace-collision", func(t *testing.T) {
		inst := f.scaffold(t, "pack-collision", "Pack Collision")
		dir := filepath.Join(inst, "program", "packs", "fixture-pack", "templates")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		const schema = `namespace = "CMP"
home = "fixture-home"
filename_pattern = "CMP-{NNN}-{slug}.md"
stage_scoped = false
digits = 3
required_front_matter = ["id"]
required_sections = ["Body"]
`
		if err := os.WriteFile(filepath.Join(dir, "collide.schema.toml"), []byte(schema), 0o644); err != nil {
			t.Fatal(err)
		}
		f.check(t, inst, 1)
	})

	t.Run("second-opinion-happy", func(t *testing.T) {
		inst := f.seedSecondOpinion(t, "second-opinion-happy")
		f.check(t, inst, 0)
	})

	t.Run("second-opinion-second-report", func(t *testing.T) {
		inst := f.seedSecondOpinion(t, "second-opinion-second-report")
		f.runOK(t, "new", "model-report", "Model B", "--dir", inst)
		ms401OverwriteOne(t, filepath.Join(inst, "reviews", "reports", "RPT-002-*.md"), ms401SecondOpinionRPT2)
		f.check(t, inst, 1)
	})

	t.Run("second-opinion-hash-nibble-flip", func(t *testing.T) {
		inst := f.seedSecondOpinion(t, "second-opinion-hash-nibble-flip")
		rpt := ms401MatchOne(t, filepath.Join(inst, "reviews", "reports", "RPT-001-*.md"))
		ms401ReplaceOne(t, rpt, appendixBHash, "fc87bfc2afd545807ca87b5c29cae8e77262cb3c746fc63e4539d8daeb2a77de")
		stderr := f.check(t, inst, 1)
		if !strings.Contains(stderr, appendixBHash) {
			t.Fatalf("stderr must name expected hash %s: %q", appendixBHash, stderr)
		}
	})

	t.Run("second-opinion-opt-in-false", func(t *testing.T) {
		inst := f.seedSecondOpinion(t, "second-opinion-opt-in-false")
		cmp := ms401MatchOne(t, filepath.Join(inst, "reviews", "commissioning", "CMP-001-*.md"))
		ms401ReplaceOne(t, cmp, "opt_in = true", "opt_in = false")
		f.check(t, inst, 1)
	})

	t.Run("second-opinion-standard-cost", func(t *testing.T) {
		inst := f.seedSecondOpinion(t, "second-opinion-standard-cost")
		cmp := ms401MatchOne(t, filepath.Join(inst, "reviews", "commissioning", "CMP-001-*.md"))
		ms401ReplaceOne(t, cmp, `cost_class = "cheap"`, `cost_class = "standard"`)
		f.check(t, inst, 1)
	})

	t.Run("cursor-council-happy", func(t *testing.T) {
		inst := f.seedCouncil(t, "cursor-council-happy", "Council Cursor", "cursor")
		f.check(t, inst, 0)
	})

	t.Run("cursor-council-missing-reconciliation", func(t *testing.T) {
		inst := f.seedCouncil(t, "cursor-council-missing-reconciliation", "Council Cursor", "cursor")
		ms401RemoveOne(t, filepath.Join(inst, "reviews", "reconciliations", "RCL-001-*.md"))
		f.check(t, inst, 1)
	})

	t.Run("cursor-council-missing-second-report", func(t *testing.T) {
		inst := f.seedCouncil(t, "cursor-council-missing-second-report", "Council Cursor", "cursor")
		ms401RemoveOne(t, filepath.Join(inst, "reviews", "reports", "RPT-002-*.md"))
		f.check(t, inst, 1)
	})

	t.Run("cursor-council-missing-retained-dissent", func(t *testing.T) {
		inst := f.seedCouncil(t, "cursor-council-missing-retained-dissent", "Council Cursor", "cursor")
		rcl := ms401MatchOne(t, filepath.Join(inst, "reviews", "reconciliations", "RCL-001-*.md"))
		ms401ReplaceOne(t, rcl, "SEED-DISSENT", "none")
		f.check(t, inst, 1)
	})

	t.Run("cursor-council-cheap-cost", func(t *testing.T) {
		inst := f.seedCouncil(t, "cursor-council-cheap-cost", "Council Cursor", "cursor")
		cmp := ms401MatchOne(t, filepath.Join(inst, "reviews", "commissioning", "CMP-001-*.md"))
		ms401ReplaceOne(t, cmp, `cost_class = "standard"`, `cost_class = "cheap"`)
		f.check(t, inst, 1)
	})

	t.Run("manual-floor-happy", func(t *testing.T) {
		inst := f.seedCouncil(t, "manual-floor-happy", "Council Manual", "manual")
		f.check(t, inst, 0)
	})

	t.Run("lone-commissioning", func(t *testing.T) {
		inst := f.scaffold(t, "lone-commissioning", "Lone Commissioning")
		f.runOK(t, "new", "commissioning", "SQLite legal default", "--dir", inst)
		f.check(t, inst, 0)
	})

	t.Run("spark-pack-present-zero-commissioning", func(t *testing.T) {
		inst := f.scaffold(t, "spark-zero-commissioning", "Spark Zero Commissioning")
		if info, err := os.Stat(filepath.Join(inst, "program", "packs", "council")); err != nil || !info.IsDir() {
			t.Fatalf("council pack missing: %v", err)
		}
		if matches, err := filepath.Glob(filepath.Join(inst, "reviews", "commissioning", "CMP-*.md")); err != nil {
			t.Fatal(err)
		} else if len(matches) != 0 {
			t.Fatalf("want zero commissioning files, got %v", matches)
		}
		f.check(t, inst, 0)
	})

	t.Run("fresh-pack-deleted-no-reviews", func(t *testing.T) {
		inst := f.scaffold(t, "fresh-pack-deleted", "Fresh Pack Deleted")
		ms401RemoveCouncilPack(t, inst)
		ms401AssertAbsent(t, filepath.Join(inst, "reviews"))
		f.check(t, inst, 0)
	})

	t.Run("uppercase-prompt-sha256", func(t *testing.T) {
		inst := f.seedSecondOpinion(t, "uppercase-prompt-sha256")
		rpt := ms401MatchOne(t, filepath.Join(inst, "reviews", "reports", "RPT-001-*.md"))
		ms401ReplaceOne(t, rpt, appendixBHash, strings.ToUpper(appendixBHash))
		f.check(t, inst, 1)
	})

	t.Run("report-rung-mismatch", func(t *testing.T) {
		inst := f.seedSecondOpinion(t, "report-rung-mismatch")
		rpt := ms401MatchOne(t, filepath.Join(inst, "reviews", "reports", "RPT-001-*.md"))
		ms401ReplaceOne(t, rpt, `rung = "second-opinion"`, `rung = "council"`)
		f.check(t, inst, 1)
	})

	t.Run("reconciliation-second-opinion-rung", func(t *testing.T) {
		inst := f.seedCouncil(t, "reconciliation-second-opinion-rung", "Council Cursor", "cursor")
		rcl := ms401MatchOne(t, filepath.Join(inst, "reviews", "reconciliations", "RCL-001-*.md"))
		ms401ReplaceOne(t, rcl, `rung = "council"`, `rung = "second-opinion"`)
		f.check(t, inst, 1)
	})

	t.Run("council-command-is-unknown", func(t *testing.T) {
		code, _, stderr := f.run(t, "council")
		if code == 0 {
			t.Fatalf("mycelium council must be unknown: stderr=%q", stderr)
		}
		if !strings.Contains(stderr, "unknown command") {
			t.Fatalf("stderr must identify unknown command: %q", stderr)
		}
	})

	t.Run("gh-never-invoked", func(t *testing.T) {
		f.assertGhNever(t)
		assertNoHomeTouch(t, f.home)
	})
}

func (f ms401Fixture) run(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	code, stdout, stderr := clitest.Run(t, f.bin, f.work, f.env, args...)
	assertNoGh(t, stderr)
	f.assertGhNever(t)
	assertNoHomeTouch(t, f.home)
	return code, stdout, stderr
}

func (f ms401Fixture) runOK(t *testing.T, args ...string) {
	t.Helper()
	code, _, stderr := f.run(t, args...)
	if code != 0 {
		t.Fatalf("%s: exit %d stderr=%q", strings.Join(args, " "), code, stderr)
	}
}

func (f ms401Fixture) scaffold(t *testing.T, name, title string) string {
	t.Helper()
	inst := filepath.Join(f.work, name)
	ms401AssertAbsent(t, inst)
	f.runOK(t, "new", "idea", title, "--offline", "--dir", inst)
	return inst
}

func (f ms401Fixture) check(t *testing.T, inst string, want int) string {
	t.Helper()
	code, _, stderr := f.run(t, "check", "--dir", inst)
	if code != want {
		t.Fatalf("check %s: exit %d want %d stderr=%q", inst, code, want, stderr)
	}
	return stderr
}

func (f ms401Fixture) seedSecondOpinion(t *testing.T, name string) string {
	t.Helper()
	inst := f.scaffold(t, name, "SO Fixture")
	f.runOK(t, "new", "commissioning", "SQLite second opinion", "--dir", inst)
	f.runOK(t, "new", "model-report", "Model A", "--dir", inst)
	ms401OverwriteOne(t, filepath.Join(inst, "reviews", "commissioning", "CMP-001-*.md"), ms401AppendixBCMP)
	ms401OverwriteOne(t, filepath.Join(inst, "reviews", "reports", "RPT-001-*.md"), ms401AppendixBRPT)
	return inst
}

func (f ms401Fixture) seedCouncil(t *testing.T, name, title, adapter string) string {
	t.Helper()
	inst := f.scaffold(t, name, title)
	f.runOK(t, "new", "commissioning", "SQLite council", "--dir", inst)
	f.runOK(t, "new", "model-report", "Model A", "--dir", inst)
	f.runOK(t, "new", "model-report", "Model B", "--dir", inst)
	f.runOK(t, "new", "reconciliation", "SQLite council reconciliation", "--dir", inst)

	cmp := ms401AppendixCCMP
	rpt1 := ms401AppendixCRPT1
	rpt2 := ms401AppendixCRPT2
	if adapter == "manual" {
		cmp = strings.ReplaceAll(cmp, `adapter = "cursor"`, `adapter = "manual"`)
		rpt1 = strings.ReplaceAll(rpt1, `adapter = "cursor"`, `adapter = "manual"`)
		rpt2 = strings.ReplaceAll(rpt2, `adapter = "cursor"`, `adapter = "manual"`)
	} else if adapter != "cursor" {
		t.Fatalf("unsupported council adapter %q", adapter)
	}

	ms401OverwriteOne(t, filepath.Join(inst, "reviews", "commissioning", "CMP-001-*.md"), cmp)
	ms401OverwriteOne(t, filepath.Join(inst, "reviews", "reports", "RPT-001-*.md"), rpt1)
	ms401OverwriteOne(t, filepath.Join(inst, "reviews", "reports", "RPT-002-*.md"), rpt2)
	ms401OverwriteOne(t, filepath.Join(inst, "reviews", "reconciliations", "RCL-001-*.md"), ms401AppendixCRCL)
	return inst
}

func (f ms401Fixture) assertGhNever(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(f.ghMarker); err == nil {
		t.Fatal("gh was invoked")
	} else if !os.IsNotExist(err) {
		t.Fatalf("stat gh marker: %v", err)
	}
}

func ms401MatchOne(t *testing.T, pattern string) string {
	t.Helper()
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatalf("want one match for %q, got %v", pattern, matches)
	}
	return matches[0]
}

func ms401OverwriteOne(t *testing.T, pattern, body string) {
	t.Helper()
	if err := os.WriteFile(ms401MatchOne(t, pattern), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func ms401RemoveOne(t *testing.T, pattern string) {
	t.Helper()
	if err := os.Remove(ms401MatchOne(t, pattern)); err != nil {
		t.Fatal(err)
	}
}

func ms401ReplaceOne(t *testing.T, path, old, replacement string) {
	t.Helper()
	body := string(readFile(t, path))
	if strings.Count(body, old) != 1 {
		t.Fatalf("want one %q in %s", old, path)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(body, old, replacement, 1)), 0o644); err != nil {
		t.Fatal(err)
	}
}

func ms401RemoveCouncilPack(t *testing.T, inst string) {
	t.Helper()
	if err := os.RemoveAll(filepath.Join(inst, "program", "packs", "council")); err != nil {
		t.Fatal(err)
	}
}

func ms401AssertAbsent(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("%s must not exist: %v", path, err)
	}
}
