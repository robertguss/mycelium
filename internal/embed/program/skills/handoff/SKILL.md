---
name: handoff
description: >
  Write the Mycelium implementation packet and set state to handed-off.
  Use when the idea is clarified and the human is ready to implement.
  Legal only from clarified. Does not reopen locked decisions.
---

# Handoff

Legal only from `clarified`.

## Steps

1. Confirm `state` is `clarified`. If not, run `clarify` first.
2. `mycelium handoff`
3. `mycelium check`

Done: `handoff/PACKET.md` exists, `state` is `handed-off`, `check` is green.

If `handoff/PACKET.md` already exists, run `mycelium state handed-off`
instead of `handoff` again.

The implementer receives **only** `handoff/`. See
`program/reference/implementation-systems.md`.

Do not `git commit` unless the human asks.
