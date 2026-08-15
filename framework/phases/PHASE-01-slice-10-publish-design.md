# Slice 10 — `internal/publish` architect Phase B

Binding SoT: brief §8.2, §8.6, §10 publish split, §13 Slice 10, §14(b), §15; DEC-010.
Existing seams: `execrun.Runner`, `op.Begin`/`Stage`/`Commit`, `manifest.GithubRepo`+`Encode`, `teach.Write`, scaffold Offline/Publish, cli publish stub, journal `op: publish`.

---

## Design A — phase machine (converge)

Organizing structure: **state machine** (`Phase` enum + observe → advance).

### 1. Caller usage

```go
// cli: mycelium publish
root, _ := check.FindRoot(start)
code := publish.Run(publish.Options{
  Root: root, Mode: publish.ModeRequire, Argv: argv,
}, publish.Deps{Clock: c, Runner: r, Stdout: out, Stderr: err})

// scaffold after git init (default auto / --publish)
res := publish.RunResult(...) // or Run returning (Result, int)
switch res.Kind {
case publish.KindPublished:
  fmt.Fprintf(out, "published: %s (topic: idea)\n", res.OwnerRepo)
case publish.KindSkippedUnauth:
  fmt.Fprintln(out, "publish: mycelium publish")
case publish.KindAlready:
  fmt.Fprintf(out, "already published: %s\n", res.OwnerRepo)
}
```

### 2. Core types

```go
type Phase int // NeedAuth, NeedCreate, NeedRemote, NeedTopic, NeedManifest, Done

type Mode int // ModeRequire (--publish | mycelium publish), ModeAuto (default new idea)

type Options struct {
  Root, RepoName string // RepoName empty → manifest.Slug; tests pass mycelium-ms101-<unix>
  Mode Mode
  Argv []string
  Offline bool // if true, Run refuses before any gh
}

type Result struct {
  Kind      Kind // Published | Already | SkippedUnauth | Failed
  OwnerRepo string
  CreatedName string // name after successful create (for cleanup/tests)
}

type machine struct {
  phase Phase
  owner, name, url string
  remoteOK, topicOK, manifestOK bool
}
```

### 3. Public signatures

```go
func Run(opts Options, deps Deps) (Result, int)
func FixtureName(now time.Time) string // mycelium-ms101-<unix>
func AllowDelete(name string) bool     // ^mycelium-ms101-[0-9]+$
```

Unexported: `observe(runner) machine`, `advance(m) (machine, error)`, `ghAuth/Create/Remote/Topic/Delete`.

### 4. Module map

```text
internal/publish/
  publish.go   // Run, Options, Result, Mode, Kind
  machine.go   // Phase, observe, advance loop
  gh.go        // auth/create/remote/topic/delete via execrun
  name.go      // FixtureName, AllowDelete, resolveRepoName
  publish_test.go
  publish_github_test.go //go:build github_integration
```

### 5. Idempotency / cleanup / journal

- `observe` sets flags from `gh` + `git remote` + `manifest.GithubRepo`.
- Loop `advance` until `Done`. Same end state from any partial start (brief idempotent cases).
- After create: `sess.SetOriginalID(createdName)` before further gh; resume reuses `OriginalID` as name (no second create).
- Failure after create: if `AllowDelete(name)` then `gh repo delete --yes`; else teach + print URL.
- Manifest/log via `op` after gh success only.

### 6. Pros / cons

Pros: idempotency is the model; crash points map to phases; brief cases are table rows.
Cons: phase enum risks temporal decomposition if files split by step; richer types than callers need; easy to leak Phase on the public surface (shallow module).

Red-flag screen: keep Phase unexported or reject A as public API. Machine is fine as private.

---

## Design B — linear Run + Outcome DU

Organizing structure: **discriminated Outcome** + one linear `Run`; small private gh helpers.

### 1. Caller usage

```go
// cli
res, code := publish.Run(publish.Options{
  Root: root, Auth: publish.AuthRequired, Argv: argv,
}, deps)
// map res to stdout inside Run (Already / Published messages)

// scaffold
res, code := publish.Run(publish.Options{
  Root: target, Auth: publish.AuthOptional, RepoName: "", Argv: argv,
}, deps)
if res.Status == publish.StatusSkipped {
  // keep "publish: mycelium publish" line
} else if res.Status == publish.StatusOK {
  // replace with published: owner/name
}
```

### 2. Core types

```go
type AuthPolicy int // AuthRequired, AuthOptional

type Status int // StatusOK, StatusAlready, StatusSkipped, StatusErr

type Options struct {
  Root, RepoName string
  Auth AuthPolicy
  Offline bool
  Argv []string
}

type Outcome struct {
  Status    Status
  OwnerRepo string // owner/name when known
  RepoName  string // short name (journal / cleanup)
}

type Deps struct {
  Clock clock.Clock
  Runner execrun.Runner
  Stdout, Stderr io.Writer
}
```

### 3. Public signatures

```go
func Run(opts Options, deps Deps) (Outcome, int)
func FixtureRepoName(unix int64) string
func CanAutoDelete(repoName string) bool
```

### 4. Module map

```text
internal/publish/
  publish.go  // Run (linear), Options, Outcome, AuthPolicy, Status
  gh.go       // authStatus, createPrivate, addOrigin, addIdeaTopic, deleteRepo
  policy.go   // FixtureRepoName, CanAutoDelete, description truncate
  *_test.go
```

### 5. Idempotency / cleanup / journal

- Early exits inside `Run`: already published → StatusAlready; optional+unauth → StatusSkipped; offline → teach refuse.
- Partial repair inline: remote exists but no topic → add topic + manifest if needed; no remote → create path.
- `op.Begin(Op:"publish")`; on create success `SetOriginalID(repoName)`; stage manifest+log; Commit.
- Cleanup: `CanAutoDelete` gates delete; user slugs never auto-delete.

### 6. Pros / cons

Pros: deep module (one Run); matches tiercmd/scaffold shape; smallest public surface; easy to wire.
Cons: idempotent branches live as conditionals in one function unless an observe helper is added; less explicit crash taxonomy than a phase enum.

Red-flag screen: pass if helpers stay private and Run does not re-export gh wire shapes. Avoid pass-through wrappers.

---

## Synthesis (chosen)

**Base: Design B.** Graft from A: a private `world` observe struct (not a public Phase API).

### Why

- Callers need one call and an Outcome. They do not drive phases. Prefer interface depth (B) over exported lifecycle (A).
- Existing packages return exit codes from `Run` with teach errors. B matches that seam.
- Idempotency still earns a structure: private `world` (remote URL, topic present, manifest set, created name from journal) is observe-first converge without temporal module split.
- Laziness: A’s public machine adds reader load without hiding more policy.

### Deviations from brief

None.

### Implementer package (contract)

```go
package publish

// AuthRequired: mycelium publish and new idea --publish (fail if unauthenticated).
// AuthOptional: default new idea (skip publish if gh missing/unauth; local success).
type AuthPolicy int
const ( AuthRequired AuthPolicy = iota; AuthOptional )

type Status int
const (
  StatusOK Status = iota // created or repaired to published
  StatusAlready          // remote+topic+github_repo already set
  StatusSkipped          // AuthOptional and no usable gh auth
  StatusErr
)

type Options struct {
  Root     string
  RepoName string // empty → manifest.Slug; credentialed tests pass FixtureRepoName
  Auth     AuthPolicy
  Offline  bool // true → refuse before LookPath/Run of gh (defense in depth)
  Argv     []string
}

type Outcome struct {
  Status    Status
  OwnerRepo string // "owner/name"
  RepoName  string // short name; journal OriginalID after create
}

type Deps struct {
  Clock            clock.Clock
  Runner           execrun.Runner
  Stdout, Stderr   io.Writer
}

func Run(opts Options, deps Deps) (Outcome, int)
func FixtureRepoName(unixSec int64) string // mycelium-ms101-<unix>
func CanAutoDelete(name string) bool       // ^mycelium-ms101-[0-9]+$

// private world: AuthOK, Owner, Name, OriginURL, TopicOK, ManifestRepo, CreatedName
// Run:
//   1. refuse Offline; Find/load manifest+log; resolve Name (OriginalID | RepoName | Slug)
//  2. op.Begin publish; observe world via gh/git
//  3. if complete → Already, Commit no-op or skip op, stdout already published
//  4. AuthOptional && !AuthOK → Skipped (no teach fail)
//  5. AuthRequired && !AuthOK → teach fail
//  6. create if needed (private, desc idea: <name>≤80); SetOriginalID; remote add; topic idea
//  7. on failure after create: CanAutoDelete → gh repo delete --yes; else teach+URL
//  8. stage github_repo + updated_date + log publish line; Commit
//  9. StatusOK stdout published: owner/name (topic: idea)
```

### Module map (final)

```text
internal/publish/publish.go   // Run, types
internal/publish/world.go     // observe + converge helpers (unexported)
internal/publish/gh.go        // execrun wrappers only
internal/publish/policy.go    // FixtureRepoName, CanAutoDelete, desc
internal/cli                  // wire publish; drop stub
internal/scaffold             // drop --publish stub; call Run after git init
```

### Wiring notes

- `MYCELIUM_OFFLINE=1` → scaffold sets Offline; never calls Run with network, or calls with Offline true.
- `--offline`+`--publish` stays scaffold refuse (already landed).
- Credentialed tests: `RepoName: FixtureRepoName(...)`, `t.Cleanup` delete via CanAutoDelete path; build tag `github_integration`.
- Journal `original_id` holds created short name so resume does not create a second repo.

### Next implementation step

Add `internal/publish` with `Run`/`Outcome` stubs and hermetic unit tests for policy + Offline refuse, then fill gh/world against `execrun.Fake`.
