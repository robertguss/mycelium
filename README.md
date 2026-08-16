# mycelium

Mycelium is a CLI and a set of conventions for thinking **with an LLM**.

You and an agent work inside a git repository dedicated to one idea. The agent
holds positions instead of nodding along. Disagreement is written down, not
lost in chat. When you walk away for weeks, you can come back without rereading
a transcript.

The CLI does not think. It creates the idea repo, writes structured artifacts,
checks that the record is intact, and tracks where the idea sits in its life.
The thinking happens in your agent.

## Why it exists

Chat is a terrible system of record. Sessions die. Assumptions evaporate.
You cannot remember why you decided something, or what would change your mind.

Mycelium gives every idea a durable home: terms, questions, positions,
decisions, dissent, evidence, and time, all in git.

It stops before a coding backlog. The last artifact is a **handoff packet**
that an implementation agent can pick up and build from.

## Install

Release binaries ship for **linux-amd64** and **darwin-arm64**. See
[`docs/install.md`](docs/install.md) for the one-liner and for building from
source. Then:

```bash
mycelium version
```

## Five-minute first idea

1. Install the CLI.
2. Scaffold a local idea (no GitHub required):

   ```bash
   mycelium new idea "garden lighting" --offline
   cd garden-lighting
   ```

3. Open that folder in Cursor, Claude Code, Grok, or any agent that can read
   files and run commands.
4. Tell the agent what you want, in plain language. For example:
   *“This is a new mycelium idea. Capture the spark.”*
5. The agent should load the shipped skills under `.agents/skills/`, interview
   you, write the first artifact, and run `mycelium check`.
6. You commit. The CLI never does.

That is the normal way to work. Typing commands yourself is supported; it is
not the product.

## How an idea lives

```text
spark → exploring ⇄ simmering → clarified → handed-off
any → archived
```

| State | Meaning |
| --- | --- |
| `spark` | Repo exists. Framing may not. |
| `exploring` | You are working. |
| `simmering` | Parked on purpose. Needs a revisit date or event. |
| `clarified` | Destination reached. A handoff packet can be built. |
| `handed-off` | Packet delivered to an implementation system. |
| `archived` | Dead or absorbed. Files stay. |

`wake` is not a seventh state. It is the ritual that takes a simmering idea
back to exploring and writes a re-entry brief.

How long you stay in each state depends on the session, not on the tool. A
thirty-minute capture and a multi-week research program use the same spine.

## Scenarios

| You want | Do this |
| --- | --- |
| Capture a thought tonight | New idea, open it in your agent, ask it to capture the spark. |
| Work a live question | Open the idea, say what you are deciding. |
| Park it | Ask the agent to simmer it with a revisit date. |
| Come back after a gap | Ask the agent to wake the idea. Read the brief, not the raw log. |
| Survey every idea | `mycelium status --all` (or ask the agent for the portfolio). |
| Run a full research program | Raise the rigor tier. Follow discovery → blueprint → charter → research. |
| Hand it to implementation | Clarify, then handoff. |

The full walkthrough is [`docs/user-guide.md`](docs/user-guide.md).

## This repository vs an idea repository

This repository **builds** the `mycelium` binary. It is not an idea.

An idea repository has `mycelium.toml` at its root. Create one with
`mycelium new idea`. Agent rules for an idea live in that repo’s `AGENTS.md`
(copied from `program/skeleton/AGENTS.md`), not in this file.

`program/` is the methodology library. The CLI embeds it and copies it into
every new idea.

## Developing the CLI

Go **1.26**. `CGO_ENABLED=0`.

```bash
CGO_ENABLED=0 go build -o mycelium ./cmd/mycelium
CGO_ENABLED=0 go test ./...
```

## Read next

- [`docs/install.md`](docs/install.md) — install the binary or build from source
- [`docs/user-guide.md`](docs/user-guide.md) — install through every workflow and command
