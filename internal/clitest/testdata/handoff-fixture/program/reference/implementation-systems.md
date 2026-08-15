# Implementation systems

Portable mapping from the handoff packet to the isolated implementer.
Default system is **pstack/poteto-mode**. `manual` is the floor and always legal.

## Packet → pstack / poteto

| Packet section | pstack / poteto | Isolated implementer does |
| --- | --- | --- |
| Framing | why/how context; poteto-mode constraints | Read first. Do not reopen. |
| Locked decisions | decided list + DEC copies in `handoff/decisions/` | Treat as Accepted. Do not restage. |
| Glossary | shared language (`handoff/glossary.md`) | Use these terms. Challenge drift only inside the packet. |
| Open questions | agreement states; poteto candor ("no is an acceptable answer") | Honor `open` / `aligned` / `agree-to-disagree`. Do not silently flip. |
| Evidence summary | citations (`handoff/evidence/`) | Cite; do not invent EVDs. |
| Implementation playbooks | how/ vertical slices for the bounded target | Execute the slices. Default system is pstack/poteto-mode. |
| Implementation system | default `pstack/poteto`; `manual` is the floor | If pstack is unavailable, `manual` still satisfies the contract. |
| Time budget | required; fixture uses `30m` | Stop at the budget. Do not expand scope. |
| Acceptance | executable tests in `handoff/acceptance/` | Run them. Green = the bounded target is done. |

## Isolation

1. The isolated implementer receives **ONLY** `handoff/`.
2. No chat history.
3. No instance source beyond the packet.
4. Do not fetch the rest of the instance. Do not reopen locked decisions. Do not treat v1 attachment manifests (`program/contracts/handoffs.md`) as the packet.

The packet is `handoff/PACKET.md`. Generate it with `mycelium handoff` (legal only from `clarified`).
