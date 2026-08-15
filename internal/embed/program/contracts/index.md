# Instance Index — index.md (Mycelium 2.0)

Required file at all tiers. Structure only (DEC-005). Check validates presence
of required headings; body is not graded.

## Required structure

```text
# <idea_name>
## State
## Artifacts
## Log tail
## Wake
```

| Heading | Generator writes | Check |
| --- | --- | --- |
| `# <idea_name>` | Concrete `idea_name` from the manifest | File has an H1. H1 text need not equal `idea_name`. |
| `## State` | `state`, `tier`, `revisit`, `github_repo` as concrete values | H2 present |
| `## Artifacts` | Counts per registered NS (files matching the filename pattern). Missing home = 0. | H2 present |
| `## Log tail` | Last 20 parseable log lines (not headings). Fewer than 20 → all of them. | H2 present |
| `## Wake` | `briefs/LATEST.md` if that file exists, else `none` | H2 present |

Extra H2s are allowed. Body prose is not graded. Tokens are not required.

## Emission rules

| Rule | Binding |
| --- | --- |
| New scaffolds | Generator writes concrete `index.md` (not tokens). |
| Existing PHASE-01 instances | May lack `index.md`. `mycelium index [--dir PATH]` rebuilds. `check` **FAILS** if missing. |
| Regeneration | Every mutating command regenerates it: `state`, `wake`, `new <type>`, `tier`, `publish`, `index`. Scaffold writes it at birth. |
| Skills | `mycelium index` does **not** emit skills. |
| Log | `mycelium index` does **not** append a log line. |

Success stdout:

```text
wrote index.md
```

## Allowed top-level additions

Add to the PHASE-01 always-allowed set:

```text
index.md
briefs/
```

`briefs/` is not an ID-to-path home. Extra files inside `briefs/` do not fail
ID-to-path.

If the log contains a `wake` op, `briefs/LATEST.md` must exist and pass the
wake H2s (`program/contracts/wake.md`).

## Tier binds

Add `index.md` to `binds` in all three `program/tiers/*.toml` files (focused,
standard, high-assurance). The bind itself lands in Slice 2; this contract
states the rule.

Do not add `briefs/` as a bind (only allowed; required only after a wake op).

## Protocol

`mycelium index` uses the operation protocol with journal `op=index`. No log
line. No skill emit. No state change.

## Package

`internal/indexmd` with `Render(instance) []byte`. Called by scaffold, `index`,
`state`, `wake`, `new <type>`, `tier`, `publish`. No network.
