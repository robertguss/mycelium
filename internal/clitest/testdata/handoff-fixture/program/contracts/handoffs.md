# Context Handoff and Attachment Manifests

> **v1 session-attachment manifests — not the PHASE-06 packet.**
> This file is the ADRP v1 attachment-manifest contract.
> The implementation packet is `program/contracts/handoff-packet.md` (`handoff/PACKET.md`).
> Do not treat attachment manifests as the handoff packet.

## Layered context

Full artifacts remain authoritative. Handoff Digests reduce size but never
silently replace source reports.

## Attachment selection rules

Every fresh session receives an explicit attachment manifest containing:

1. Governing artifacts in full
2. The current stage prompt in full
3. Direct prerequisite reports in full
4. Accepted Decision Records
5. Handoff Digests for indirectly relevant reports
6. Full indirect reports when nuance, weak evidence, or conflict is material

Chief Architect synthesis and adversarial review should receive all materially
relevant full reports unless reliable repository retrieval is available.

## Templates

- [`../templates/attachment-manifest.md`](../templates/attachment-manifest.md)
- [`../templates/launch-message.md`](../templates/launch-message.md)

## Just-in-time stage package

For every substantive stage, produce five items:

1. Canonical stage prompt
2. Repository installation task
3. Attachment manifest
4. Fresh-session launch message
5. Post-stage validation task and recommended commit message

See [`../templates/stage-package.md`](../templates/stage-package.md).
