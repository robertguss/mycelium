# Idea Lifecycle (Mycelium 2.0)

Exact vocabulary (DEC-006):

```text
spark → exploring ⇄ simmering → clarified → handed-off
any → archived
```

Birth state is always `spark`.

## Legal transitions

| From | Allowed next states |
| --- | --- |
| `spark` | `exploring`, `archived` |
| `exploring` | `simmering`, `clarified`, `archived` |
| `simmering` | `exploring`, `archived` |
| `clarified` | `handed-off`, `archived` |
| `handed-off` | `archived` |
| `archived` | *(none)* |

## PHASE-01 storage rules

PHASE-01 ships **no** state-transition command. Check still validates stored
state:

- Stored `clarified` or `handed-off` → check **FAIL** (PHASE-02 / PHASE-06 not
  shipped; those states are not reachable by a shipped command yet and must not
  appear in fixtures without their phase contracts).
- Stored `archived` is legal.
- `simmering` requires `revisit` non-empty in `mycelium.toml`.

Wake is the `simmering` → `exploring` transition (PHASE-02 ritual; not a
PHASE-01 command).
