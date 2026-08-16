---
name: simmer
description: >
  Park a Mycelium idea on purpose: exploring → simmering with a required
  revisit date or event. Use when the human wants to stop, come back later,
  or says park/simmer. Not for blocked work and not for unstateable fog.
---

# Simmer

Park because you *could* decide now and are choosing not to.

## Steps

1. Agree a revisit: `YYYY-MM-DD` or `event:<kebab>`.
2. `mycelium state simmering --revisit VALUE`
3. `mycelium check`

Done: `state` is `simmering` and `revisit` is set.

`--revisit` is required. Empty or missing is refused.

Do not `git commit` unless the human asks.
