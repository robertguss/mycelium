# Handoff Packet

This is the PHASE-06 **implementation packet** contract. The instance object is `handoff/PACKET.md`.

`program/contracts/handoffs.md` is **not** the packet. That file is the ADRP v1 session-attachment manifest contract. Do not generate, check, or emit the packet from `handoffs.md`.

## Instance tree

The packet lives at the instance root under `handoff/` (an allowed top-level path, like `briefs/`):

```text
handoff/
  PACKET.md
  decisions/
  glossary.md
  questions/
  evidence/
  playbooks/
  acceptance/
```

Never `docs/handoffs/`. Never `framework/`. Never emit this tree under those paths.

## Front matter

Required keys on `handoff/PACKET.md`. Unknown keys are allowed (append-only). Required keys must be present and valid.

| Key | Rule |
| --- | --- |
| `id` | `HO-001` (one packet per instance this phase). Allocated by the handoff generator, **not** by `mycelium new`. Do not add a `HO` namespace to `new`. |
| `date` | `YYYY-MM-DD` |
| `implementation_system` | `pstack/poteto` \| `manual`. Generator writes `pstack/poteto`. `manual` is the floor. |
| `time_budget` | Required. Regex `^[0-9]+[mh]$`. Generator writes `30m`. |

## Required H2s (order fixed)

1. Framing
2. Locked decisions
3. Glossary
4. Open questions
5. Evidence summary
6. Implementation playbooks
7. Implementation system
8. Time budget
9. Acceptance

Check validates **containers only** (DEC-005): H2s present in this order, front matter valid, listed copies exist, ID / path links resolve inside the packet. Do **not** grade Framing prose. Do **not** grade playbook quality.

## Self-contained

The packet is self-contained. ID and path links must resolve **inside** `handoff/` copies (`handoff/decisions/`, `handoff/questions/`, `handoff/evidence/`, playbooks, acceptance).

A link to `../decisions/DEC-001-*.md` (instance tree) is a **FAIL** even if that file exists outside the packet. Check **FAILS** if `PACKET.md` or any file under `handoff/playbooks/` links to an instance path outside `handoff/` that was not copied in.

## Generator copy table

`mycelium handoff` builds `handoff/` from the instance, then flips state. `HO-001` is allocated by this generator, not `mycelium new`.

| Packet path | Source |
| --- | --- |
| `PACKET.md` | Template filled from instance + defaults (`pstack/poteto`, `30m`, `HO-001`). H2 bodies are structural lists (IDs, titles, paths) — not graded prose. |
| `decisions/` | Copy every instance DEC whose `status = "Accepted"`. If none, directory exists empty and Locked decisions H2 lists `none`. |
| `glossary.md` | Copy instance `CONTEXT.md` (or emit a file whose body is `none` if CONTEXT is empty). |
| `questions/` | Copy every instance OQ, preserving `agreement`. If none, directory exists empty and Open questions H2 lists `none`. |
| `evidence/` | Copy every EVD cited from Accepted DECs or from instance `evidence/`. If none, directory exists empty and Evidence summary H2 lists `none`. Write `handoff/evidence/SUMMARY.md` listing cited IDs. |
| `playbooks/` | Copy instance `playbooks/` if present and non-empty; else emit `handoff/playbooks/PLAYBOOK.md` stub (H2s: Target, Steps, Done). |
| `acceptance/` | Copy instance `acceptance/` if present and non-empty; else emit `handoff/acceptance/README.md` stating `none` (structure pass, not a golden-impl pass). |

A clarified instance with no playbooks and no acceptance still produces a **structurally** complete packet (stubs + `none`).
