# Sparring (Mycelium 2.0)

Implements **DEC-007**. Structure only (DEC-005). Check validates containers,
never contents. No new CLI verb. Disagreement lives **ON the OQ file**. No
`DSG` namespace. No `session.md` type.

Parser package: `internal/sparring` lands in Slice 1. This contract is the
spec. Check does not call the parser until Slice 2.

## Agreement enum

```text
open | aligned | agree-to-disagree
```

Invalid `agreement` → FAIL. Teaching error names
`program/templates/question.schema.toml` and the illegal token.

| Value | Meaning | Terminal-by-convention? | Disagreement record? |
| --- | --- | --- | --- |
| `open` | Still in play. Positions container present. | No. Edit in place. | **Not** required. Extra Crux/Reasons do not fail. |
| `aligned` | Parties agree. Honorable close. | **Yes.** Do not edit back to `open`; open a new OQ. | **Not** required. |
| `agree-to-disagree` | Parties disagree and stop. Terminal and honorable. | **Yes.** Do not edit back to `open`; open a new OQ. | **Required.** |

Terminal-by-convention is skill + contract. Check does not keep agreement
history and does not fail a flip back to `open`.

## Always-required H2s (every OQ)

| Heading | Rule |
| --- | --- |
| `## Question` | Required. Body not graded. |
| `## Context` | Required. Body not graded. |
| `## Positions` | Required as a **container**. Body not graded. |
| `## Disposition` | Required. Body not graded. |

Heading match is exact and case-sensitive. Extra H2s are allowed.

## IFF `agreement == "agree-to-disagree"`

All of the following are required. Missing any → FAIL. Teaching error names
the **missing heading** and this contract.

| Heading | Nested H3s |
| --- | --- |
| `## Positions` | **Must** contain `### Human` and `### Agent` |
| `## Reasons` | **Must** contain `### Human` and `### Agent` |
| `## Crux` | **Must** contain `### Human` and `### Agent` |

H3 names are exact: `### Human` and `### Agent`. H3s must appear in the
parent H2's section body (bytes after that H2 until the next H2 or EOF).
Extra H3s in the section do not fail. H3 body may be `<!-- fill -->` or one
word — do not grade the words (DEC-005).

## IFF `agreement == "aligned"` or `open`

| Heading | Check |
| --- | --- |
| `## Positions` | Required as a container. H3 Human/Agent **not** required. |
| `## Crux` | **Not** required. |
| `## Reasons` | **Not** required. |
| Extra Crux / Reasons / H3s | Do **not** fail. |

Aligned session: **no disagreement record required**.

## Schema cannot express IFF

`question.schema.toml` `required_sections` becomes
`["Question", "Context", "Positions", "Disposition"]`. Drop `Crux` from the
always-required list (Slice 1). Do not add `Reasons`. Do not invent a schema
DSL for conditionals. IFF rules bind in check in Slice 2.

## Always-on (skill)

Sparring is always-on (skill) in any non-archived state. Check does **not**
require any OQ to exist. Spark with zero questions still passes.

## Optional DEC `## Dissent`

| Rule | Binding |
| --- | --- |
| Heading absent | Pass. Do not require Dissent on existing or new DECs. |
| Heading present | Section must contain at least one resolvable `OQ-###` or `ASM-###`. |
| Prose | Not graded. Citing only a `DEC-###` does not satisfy. |

Teaching error:

```text
mycelium: DEC-001 ## Dissent has no resolvable OQ-### or ASM-###
convention: dissent
contract: program/contracts/sparring.md
fix: cite an existing OQ-### or ASM-### in ## Dissent, or remove the heading
```

## DEC-005

Do not grade Positions / Reasons / Crux prose. Do not score substantive vs
bare. That judgment belongs to the human or an adversarial reviewer.

## Teaching-error examples

```text
mycelium: OQ-001 missing ## Crux (required when agreement=agree-to-disagree)
convention: sparring
contract: program/contracts/sparring.md
fix: add ## Crux with ### Human and ### Agent

mycelium: OQ-001 missing ### Human under ## Positions
convention: sparring
contract: program/contracts/sparring.md
fix: add ### Human and ### Agent under ## Positions

mycelium: OQ-001 missing ## Reasons (required when agreement=agree-to-disagree)
convention: sparring
contract: program/contracts/sparring.md
fix: add ## Reasons with ### Human and ### Agent

mycelium: agreement "maybe" is not open|aligned|agree-to-disagree
convention: question-front-matter
contract: program/templates/question.schema.toml
fix: set agreement to open, aligned, or agree-to-disagree
```
