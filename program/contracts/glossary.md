# Glossary (Mycelium 2.0)

`CONTEXT.md` glossary rules. Structure only (DEC-005). Skill challenges
drift and vagueness; check validates containers only.

## File existence

Do **not** add a new required-file bind. New scaffolds already emit
`CONTEXT.md`. If the file is missing, do not fail solely for that. If it
exists, apply the rules below.

## H1

If `CONTEXT.md` exists, it must contain a line that is exactly `# Glossary`.

## Empty glossary

H1 only is **legal**. Check does not require N terms.

## Terms

| Rule | Binding |
| --- | --- |
| Any `## <Term>` | That term's section **must** contain H3 `Definition`. Missing → FAIL. |
| Definition body | Not graded. `<!-- fill -->` or one word does **not** fail. |
| Content | Check does not grade definitions, challenge quality, or drift. |

Teaching error names the term heading and this contract:

```text
mycelium: CONTEXT.md term "SQLite" missing ### Definition
convention: glossary
contract: program/contracts/glossary.md
fix: add ### Definition under ## SQLite
```

## Assumption audit (skill-only)

| Rule | Binding |
| --- | --- |
| File type | **None.** Do not add an `AUDIT` type or namespace. |
| Command | Agent uses existing `mycelium new assumption`. |
| Check | Does **not** require a periodic audit file. Does **not** require N assumptions. |

Do not reimplement the assumption template or `mycelium new assumption`.
