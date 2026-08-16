---
name: clarify
description: >
  Mark a Mycelium idea clarified: destination reached, handoff packet
  buildable. Use when the human says the idea is decided enough to hand
  off, or asks to clarify. Does not write the packet.
---

# Clarify

Legal from `exploring` only.

## Steps

1. Confirm the destination is actually reached (locked decisions the
   human accepts). Do not invent “we’re done.”
2. `mycelium state clarified`
3. `mycelium check`

Done: `state` is `clarified`. Packet is **not** written.

Run `handoff` only when the human asks to implement. `clarified` without
a packet is legal.

Do not `git commit` unless the human asks.
