# Living notes — Arvo / BEAM research program

- **Repo:** [robertguss/artifact-driven-research-program](https://github.com/robertguss/artifact-driven-research-program)
- **Status:** Working notes. **Not** an accepted Blueprint or Charter.
- **Updated:** 2026-08-14
- **Owner:** Robert Guss (with Thinking Partner, Research, Engineering)

This file is the system of record for *conversation so far*. Chat is not authority. Placeholders in `docs/00-program-blueprint.md` stay placeholders until discovery framing is approved.

---

## In one paragraph

Personal lab to push the BEAM as hard as possible for an agentic coding harness (Arvo). Hypothesis is José Valim’s: the runtime *is* the framework. This program only **researches, ideates, and writes hypotheses**. A **later project** will run spikes, measure with evals, and score what gets into `arvo/`. Not a race with other harnesses. Learn from papers and other systems (including Jido), then decide.

---

## Two programs (locked)

| Phase | Where | What we do | What we do *not* do |
|-------|--------|------------|---------------------|
| **1 — this repo** | artifact-driven-research-program | Widespread research, ideation, hypotheses, a catalog of things to try | Spikes, evals, PRs into Arvo |
| **2 — later** | a new project (not stood up yet) | Implement experiments, measure, score, incorporate or drop | Invent the research agenda from scratch |

The template’s “implementation plan” in **phase 1** means: **ranked hypotheses + what we would measure + keep/drop criteria**. It does not mean code.

**Lab vs product:** the lab may try things Arvo’s daily-driver refusal list forbids (MCP, plan mode, etc.). Landing in [coding-agent-harness](https://github.com/robertguss/coding-agent-harness) `arvo/` is a **separate gate**.

---

## Hypothesis (José)

[Tweet, 14 Aug 2026](https://x.com/josevalim/status/2088186994849468659) and follow-up [“the runtime is the framework”](https://x.com/josevalim/status/2088208133487264078):

1. Hot-code swap plugins (Pi-like) without dropping state
2. Client/server as a byproduct of actors (OpenCode-like)
3. Distribution isolates **brains** (model + session) from **hands** (sandbox + tools) — “this is basically how Livebook works”

Arvo was started on that bet. Today it spends a *thin* OTP slice: singleton Session, turn Task, in-process tools, plugin load via compile + path. Focus quit still halts the VM. That is the gap, not “pick Elixir.”

---

## Three research fronts (sequence leaning locked)

| Front | Question |
|-------|----------|
| **A** | What is cutting-edge in agentic / harness engineering (mostly Python / arXiv)? |
| **B** | Which BEAM primitives and runtime gifts can we push that are *native*? |
| **C** | How do we **translate** A into a BEAM-native move (not the same scaffold on a different VM)? |

**Leaning locked — B first:** draft a BEAM primitive atlas, put papers on **Watch** in parallel, and only then **translate** (C). C remains the contribution (not a reading club, not tourist OTP) — the *sequence* is atlas → Watch → translate, not “C as day-one center.”

**Edification in scope:** RLM, GEPA, and similar work are fair game to *learn from*. Porting them into Elixir is **not** automatic and is not the default. If a paper is the same idea on another VM, it does **not** get an Elixir rewrite.

---

## Idea shelves (proposed)

Looser *intake* so Python-heavy papers survive. Not a looser *graduation* bar.

| Shelf | Meaning |
|-------|---------|
| **Watch** | Paper or primitive; no translation yet |
| **Translate** | We have a BEAM-shaped hypothesis (C fired) |
| **Graduate** | Ready for the experiment project: claim, what we’d measure, what “land in Arvo” means |

Most ideas should stay on Watch. That is success.

---

## BEAM primitive atlas (draft)

Working map of runtime gifts we might push. **Not accepted** — refine before Blueprint. Atlas work is **B**. Papers land on Watch while this drafts; **C** fires when a Watch item maps onto a cluster without becoming LangChain-on-BEAM.

| Cluster | Rough meaning (draft) | Seeds already in the notes |
|---------|----------------------|----------------------------|
| **Isolation** | Separate failure and trust domains (brain vs hands) without pretending Docker is optional | Hidden hands node; no secrets on hands; Docker for bash; Livebook-style nodes |
| **Liveness** | Keep the brain (Session, mailbox, JSONL) alive across client and plugin churn | Focus as disposable client; quit ≠ kill brain; plugin hot-load (`:code.load_binary` + ETS) |
| **Attention-as-topology** | Who talks to whom *is* the architecture — mailboxes, links, monitors; not HTTP-as-identity | Client/server as actor byproduct; Session cast/send only; José’s OpenCode-shaped point |
| **Concurrency policy** | Explicit choices about what runs where, backpressure, and recursive LLM calls | Turn as Task; tools aimed at hands; broker so the sandbox never sees API keys |

---

## Candidate spikes (not a backlog — evidence for later)

From Engineering + Research, synthesized. **No coding until phase 2 / a greenlight.**

1. **E2 — Focus as disposable client.** Already unlinked (`Task.start`). Real work: stop `:halt_on_focus_quit`; Session must not know a TUI module (cast/send only). Quit tile ≠ kill brain. Human IEx attach can sit beside this (dev-only, full trust).
2. **E3-thin — hidden hands node.** Port or `:peer`, `--hidden`, per-session cookie, narrow API, **no secrets on hands**. `:erpc read` only. Kill the node; Session + JSONL live. Shared cookie = unrestricted `:erpc` on the brain — do not do that.
3. **E1 — plugin hot-load.** `:code.load_binary` + ETS; Session mailbox stays. Tools aimed at hands; slash/UI on brain. Not OTP relups. Not mid-turn. Not the FFF NIF.
4. **E3-bash.** Same protocol; node in **Docker**. Hidden BEAM ≠ bash jail.
5. **IEx / RLM — two products.**
   - **Human IEx** on the **brain** (live debug, patch). Dev-only.
   - **Agent RLM** ([arXiv:2512.24601](https://arxiv.org/abs/2512.24601)): prompt as data in a REPL + recursive `llm_query`. Official impl is **Python**. Default: Port that env in Docker on hands. **IEx-as-RLM only if the context is BEAM-shaped** (ETS, Session dumps, mailboxes). “We’re Elixir so the REPL is IEx” is not a reason. Never `eval` on the Session VM. Recurse via a broker so the sandbox never sees API keys.

### Hard nos (for Arvo product; lab may still *study* them)

Horde / libcluster / Oban as the architecture; MCP in core; Jido or Alloy *as* Session; Legion/Dune as the bash isolation story; relups for plugins; attached Phoenix as default; cookie as auth; permission popups / plan / todo in the daily driver.

---

## What others spent (so we don’t redo them)

| System | What they used BEAM for | Not a harness |
|--------|-------------------------|---------------|
| **Livebook** | Brains vs Port-spawned / attached / remote nodes. Architecture paper: `runtime/standalone.ex` | Notebook, not an agent |
| **Tidewave** | IEx-for-the-agent *inside the target app* (`project_eval`). MCP. Docker is the real fence | Runtime intelligence, not Arvo |
| **Alloy** | Minimal OTP agent *loop* | Don’t let it own `start_turn` |
| **Jido** | Native Elixir agents | Study; don’t become LangChain-on-BEAM |
| **jido_harness** | Wraps TS CLIs (Claude Code, etc.) | Anti-pattern *for Arvo* |
| **Legion / Dune** | In-VM AST / eval sandbox | They say: real isolation = another BEAM |
| **Pi** | `/reload` recycles the extension process | Not BEAM hot-swap |
| **OpenCode** | HTTP client/server | José means message-passing is cheaper, not “build Hono” |

---

## Discovery status

Interview in progress (one question at a time). Do not accept Blueprint until framing is approved.

| # | Topic | State |
|---|--------|--------|
| 1 | Problem | **Locked:** personal lab; José hypothesis; learn + catalog; not compete |
| 2 | Done-enough for *this* phase | **Locked:** open-ended research; output is a collection of hypotheses; spikes are the *next* project |
| 3 | Graduation bar | **Leaning loose intake**; shelves proposed (Watch / Translate / Graduate) |
| 4 | Sequence (B / A / C) | **Leaning locked:** BEAM atlas first; papers Watch in parallel; translate after atlas |

Still to cover (later): rigor tier, tracks, paper-ingest habit, where phase-2 repo lives, success criteria for the catalog.

---

## Repos

| Repo | Role |
|------|------|
| [artifact-driven-research-program](https://github.com/robertguss/artifact-driven-research-program) | This program (still template placeholders — run `just init name="arvo-beam-harness"` when ready) |
| [coding-agent-harness](https://github.com/robertguss/coding-agent-harness) `arvo/` | Daily-driver Elixir harness. Ignore `ore/` unless we say so |
| [alexzhang13/rlm](https://github.com/alexzhang13/rlm) | Official RLM (Python REPL) |

---

## Open for Robert

1. Draft / react to the BEAM primitive atlas (four clusters).
2. `just init` this repo (name only) so it stops being `{{PROJECT_NAME}}`.
3. Drop arXiv links when you want them on the Watch shelf.

When you’re happy with this framing (atlas sequence + shelves), we fill `docs/00-program-blueprint.md` and you accept it. Not before.
