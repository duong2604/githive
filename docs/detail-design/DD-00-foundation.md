---
id: DD-00
title: Shared foundation — errors, output, modes, git runner, lock
status: Approved
req: [NFR-01, NFR-02, NFR-05]
depends_on: [BD-00, ADR-0002, ADR-0004]
owner: Advisor
last_updated: 2026-11-14
---

# DD-00 — Shared foundation

## 1. Scope

Implements [BD-00](../basic-design/BD-00-foundation.md), rules BR-01 … BR-13. It changes the existing `internal/cli` root from P02 and adds the packages every later DD builds on.

Deferred elsewhere:
- WSL `repo_on_mnt` warning: `platform.UnderMnt` is here (§3.6); the warning itself is built by target resolution (DD-03), so `internal/platform` never imports `internal/ui`.
- Interactive prompts for missing values: the `Mode` decision is here; each command asks through `ui.Prompter` (§3.3) in its own DD.
- Windows console code page: `doctor` (DD-09) reports it; no change at start-up.

### Implementation slices (one issue each)

| Slice | Contents | Depends on | Size | Test cases |
| --- | --- | --- | --- | --- |
| DD-00a | `internal/apperr`, `internal/ui` (envelope, printer, mode, prompter), `internal/platform` (term, wsl, paths), changes to `internal/cli` and `cmd/githive` | — | ~450 lines incl. tests (accepted above the ~400 guide, PO-01-D1) | TC-01 … TC-11, TC-25 … TC-27 |
| DD-00b | `internal/git` runner and version check | DD-00a | ~300 lines | TC-12 … TC-18 |
| DD-00c | `internal/lock` | DD-00a | ~250 lines | TC-19 … TC-24, TC-28 |

After DD-00a is merged, DD-00b, DD-00c and DD-03a can run in parallel; they touch different packages. DD-03b needs DD-00b (see DD-03 §1).

## 2. Packages & files

| Path | Responsibility |
| --- | --- |
| `internal/apperr/apperr.go` | `Error` type, `Code` catalogue, exit-code table |
| `internal/ui/envelope.go` | JSON envelope types and writer |
| `internal/ui/printer.go` | Human output: results to stdout, `githive: …` / `hint: …` to stderr, colour switch |
| `internal/ui/mode.go` | `Mode` (interactive, JSON, yes, verbose, color) computed from flags and terminals |
| `internal/ui/prompt.go` | `Prompter` interface; `huh`-based implementation; `NoPrompter` for non-interactive |
| `internal/cli/root.go` | Global flags, `Execute(ctx, …)`, error → exit code, JSON/text rendering of every command |
| `internal/cli/version.go` | `version` moved out of `root.go`, using the envelope |
| `cmd/githive/main.go` | `signal.NotifyContext` for interrupt; passes `ctx` to `cli.Execute` |
| `internal/git/runner.go` | `Runner`: runs git with an argument list, env, verbose trace, error mapping |
| `internal/git/version.go` | Parse `git version` output; `RequireMin` |
| `internal/git/redact.go` | Remove credentials from URLs in git messages |
| `internal/lock/lock.go` | Repository lock on `githive.lock` with holder file `githive.lock.info` |
| `internal/platform/term.go` | `IsTerminal(*os.File)` via `golang.org/x/term` |
| `internal/platform/wsl.go` | `IsWSL()`, `UnderMnt(path)` |
| `internal/platform/paths.go` | `NativeAbs(path)`: absolute, cleaned, OS-native separators (BR-11) |

New dependencies (ADR-0004): `golang.org/x/term`, `github.com/gofrs/flock`, `github.com/charmbracelet/huh`. The last two are already on the default list in CLAUDE.md §5.

## 3. Types & interfaces

### 3.1 `internal/apperr`

```go
package apperr

type Code string

const (
    Internal        Code = "internal"          // exit 1
    GitFailed       Code = "git_failed"        // exit 1
    IOFailed        Code = "io_failed"         // exit 1
    Usage           Code = "usage"             // exit 2
    InputRequired   Code = "input_required"    // exit 2
    NotConfigured   Code = "not_configured"    // exit 2
    BaseBranchUnset Code = "base_branch_unset" // exit 2
    InvalidTarget   Code = "invalid_target"    // exit 2
    ProjectNotFound Code = "project_not_found" // exit 2
    TaskNotFound    Code = "task_not_found"    // exit 2
    TaskExists      Code = "task_exists"       // exit 2
    BranchExists    Code = "branch_exists"     // exit 2
    GitMissing      Code = "git_missing"       // exit 2
    GitTooOld       Code = "git_too_old"       // exit 2
    SyncConflict    Code = "sync_conflict"     // exit 3
    Unsafe          Code = "unsafe"            // exit 4
    LockTimeout     Code = "lock_timeout"      // exit 5
    Interrupted     Code = "interrupted"       // exit 130
)

// Error is the only error type that crosses from services to the CLI.
type Error struct {
    Code    Code
    Message string // one line, English, no trailing period
    Hint    string // optional: a command or action the user can take
    Err     error  // optional cause, never shown in JSON
}

func (e *Error) Error() string // Message, plus ": " + Err.Error() when Err != nil
func (e *Error) Unwrap() error

func New(code Code, msg string) *Error
func Newf(code Code, format string, args ...any) *Error
func Wrap(code Code, err error, msg string) *Error
func (e *Error) WithHint(h string) *Error // returns e for chaining

// ExitCode maps a code to the exit code above; unknown codes → 1.
func ExitCode(c Code) int

// From converts any error: *Error is returned as is; context.Canceled → Interrupted;
// anything else → Internal wrapping err.
func From(err error) *Error
```

### 3.2 `internal/ui` — envelope and printer

```go
package ui

type Warning struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type ErrorBody struct {
    Code     apperr.Code `json:"code"`
    Message  string      `json:"message"`
    Hint     string      `json:"hint,omitempty"`
    ExitCode int         `json:"exit_code"`
}

type Envelope struct {
    OK       bool       `json:"ok"`
    Command  string     `json:"command"`        // "version", "project list"; "" for the root
    Data     any        `json:"data,omitempty"` // omitted on failure
    Error    *ErrorBody `json:"error,omitempty"`
    Warnings []Warning  `json:"warnings"`       // always an array, never null
}

// WriteJSON writes e as indented JSON (2 spaces) plus a newline.
func WriteJSON(w io.Writer, e Envelope) error

type Printer struct {
    Out, Err io.Writer
    Color    bool
}

func (p *Printer) Result(format string, args ...any)   // stdout
func (p *Printer) Warn(w Warning)                       // stderr: "warning: <message>"
func (p *Printer) Error(e *apperr.Error)                // stderr: "githive: <message>" [+ "hint: <hint>"]
func (p *Printer) Trace(format string, args ...any)     // stderr, verbose only (caller checks)
```

### 3.3 `internal/ui` — mode and prompter

```go
type Mode struct {
    Interactive bool // BR-01
    JSON        bool
    Yes         bool
    Verbose     bool
    Color       bool
}

// ModeInput is everything Mode depends on, so tests can set it directly.
type ModeInput struct {
    JSON, Yes, Verbose, NoColor bool
    StdinTTY, StdoutTTY         bool
    NoColorEnv                  bool // NO_COLOR is set (any value)
}

func ComputeMode(in ModeInput) Mode
// Interactive = StdinTTY && StdoutTTY && !JSON && !Yes
// Color       = StdoutTTY && !JSON && !NoColor && !NoColorEnv

type Prompter interface {
    Input(ctx context.Context, question, def string) (string, error)
    Select(ctx context.Context, question string, options []string) (string, error)
    Confirm(ctx context.Context, question string) (bool, error)
}

// NoPrompter returns apperr InputRequired for Input/Select, and for Confirm
// returns (true, nil) when yes is set, otherwise InputRequired.
func NoPrompter(yes bool) Prompter
func HuhPrompter() Prompter
```

### 3.4 `internal/git`

```go
package git

type Runner struct {
    Interactive bool      // false → GIT_TERMINAL_PROMPT=0, GCM_INTERACTIVE=never
    Trace       io.Writer // non-nil when --verbose
    gitPath     string    // resolved once with exec.LookPath("git")
    once        sync.Once
    initErr     error
}

// Run runs `git -C dir args...` and returns stdout with the trailing newline removed.
// It checks the git version on the first call of the process (BR-08).
func (r *Runner) Run(ctx context.Context, dir string, args ...string) (string, error)

// RunZ is Run for -z output: returns stdout split on NUL, dropping the final empty field.
func (r *Runner) RunZ(ctx context.Context, dir string, args ...string) ([]string, error)

type Version struct{ Major, Minor, Patch int; Raw string }

func ParseVersion(out string) (Version, error)
func (v Version) AtLeast(major, minor int) bool

const MinMajor, MinMinor = 2, 38
```

### 3.5 `internal/lock`

```go
package lock

type Holder struct {
    Command string    `json:"command"` // e.g. "githive spawn kirindo@task-18"
    PID     int       `json:"pid"`
    Host    string    `json:"host"`
    Started time.Time `json:"started"` // RFC 3339, UTC
}

type Lock struct{ /* flock + paths */ }

type Options struct {
    Timeout time.Duration    // from Timeout(); default 30 s
    OnWait  func(h *Holder)  // called once, after 1 s of waiting; h may be nil
    Self    Holder
}

// Acquire locks <commonDir>/githive.lock and writes <commonDir>/githive.lock.info.
func Acquire(ctx context.Context, commonDir string, opt Options) (*Lock, error)

// Release removes githive.lock.info, then unlocks. Safe to call twice.
func (l *Lock) Release() error

// Timeout reads GITHIVE_LOCK_TIMEOUT (whole seconds, 1–3600); unset → 30 s.
// Any other value → apperr Usage.
func Timeout(getenv func(string) string) (time.Duration, error)

// With acquires, runs fn, and always releases. Errors are returned through apperr.From,
// so a cancelled ctx (in Acquire or from fn) becomes interrupted (BR-12).
func With(ctx context.Context, commonDir string, opt Options, fn func() error) error
```

### 3.6 `internal/platform`

```go
func IsTerminal(f *os.File) bool        // golang.org/x/term
func IsWSL() bool                       // WSL_DISTRO_NAME set, or /proc/sys/kernel/osrelease contains "microsoft" (case-insensitive)
func UnderMnt(path string) bool         // IsWSL() && path matches ^/mnt/[a-z]/
func NativeAbs(path string) (string, error) // filepath.Abs + filepath.Clean; used for every path printed or put in JSON (BR-11)
```

`internal/platform` imports only the standard library and `golang.org/x/term`; it never imports `internal/ui` or `internal/apperr`.

## 4. Config keys / JSON output schema

- **Environment:** `GITHIVE_LOCK_TIMEOUT` (seconds, 1–3600). `NO_COLOR` (any value disables colour).
- **No config file keys** are introduced by this DD.
- **Envelope:** §3.2, for every command. Field order in the output follows the struct order.
- **`version --json` data:**

```json
{
  "ok": true,
  "command": "version",
  "data": { "version": "0.0.0-dev", "commit": "none", "date": "unknown" },
  "warnings": []
}
```

- **Lock info file** `githive.lock.info` (JSON, one object, written atomically):

```json
{ "command": "githive spawn kirindo@task-18", "pid": 4312, "host": "DESKTOP-1", "started": "2026-11-13T09:15:02Z" }
```

## 5. Algorithm

### 5.1 `cli.Execute(ctx, info, args, stdin, stdout, stderr) int`

1. Build the cobra tree. Global persistent flags: `--json`, `--yes/-y`, `--verbose/-v`, `--no-color`. `SilenceUsage` and `SilenceErrors` stay on.
2. Flag and argument errors (`SetFlagErrorFunc`, `Args` validators, unknown command) become `apperr.Usage`.
3. In `PersistentPreRunE`: compute `Mode` from the flags, `platform.IsTerminal(os.Stdin)`, `platform.IsTerminal(os.Stdout)` and `NO_COLOR`; store mode, printer and prompter in the command's context.
4. Each command returns `(data any, warnings []ui.Warning, err error)` through a small adapter (`run(func(ctx, deps) (any, []ui.Warning, error))`), so rendering happens in one place.
5. Rendering:
   - JSON mode: write one envelope to stdout (success or failure). Nothing goes to stderr except `--verbose` traces.
   - Text mode: the command prints its own result through `Printer.Result`; warnings go through `Printer.Warn`; an error through `Printer.Error`.
6. If `ctx.Err() != nil` when the command returns an error, the error becomes `apperr.Interrupted`.
7. Return `apperr.ExitCode(err.Code)`, or 0.

`cmd/githive/main.go`: `ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)`, adding `syscall.SIGTERM` on non-Windows builds through a build-tagged helper; `defer stop()`; `os.Exit(cli.Execute(ctx, …))`.

Running `githive` with no arguments still prints help (P01d-D3).

### 5.2 `git.Runner.Run`

1. First call: `exec.LookPath("git")`. Not found → `GitMissing` with hint `install git 2.38 or newer: https://git-scm.com/downloads`. Then run `git version` (with the env below) and `ParseVersion`. Below 2.38 → `GitTooOld`, message `git 2.38 or newer is required, found <Raw>`.
2. Build the command: `exec.CommandContext(ctx, gitPath, append([]string{"-C", dir}, args...)...)`.
3. Environment: `os.Environ()` plus `LC_ALL=C`; when `!Interactive`, also `GIT_TERMINAL_PROMPT=0` and `GCM_INTERACTIVE=never` (Git Credential Manager on Windows). Later entries override earlier ones.
4. Capture stdout and stderr into separate buffers. Stdin is nil, except in interactive mode, where it is the process stdin so git can ask for credentials.
5. When `Trace != nil`, after the call write `git -C <dir> <args joined by space> (<ms> ms)` to Trace. Arguments are shown as passed, after `redact`.
6. Errors:
   - `ctx.Err() != nil` → `Interrupted`.
   - Non-zero exit → `GitFailed`, message `git <first arg> failed: <first line of redacted stderr>`. The whole redacted stderr goes into `Err`.
7. Return stdout with one trailing `\n` removed.

`ParseVersion`: regex `^git version (\d+)\.(\d+)(?:\.(\d+))?`. It accepts `git version 2.43.0`, `git version 2.47.1.windows.1` and `git version 2.39.3 (Apple Git-146)`.

`redact`: replace `://<anything without / @ or whitespace>@` with `://***@` (removes `user:token@`).

### 5.3 `lock.Acquire`

1. `path := filepath.Join(commonDir, "githive.lock")`, `info := path + ".info"`.
2. `f := flock.New(path)`. Make a context with deadline `now + Timeout`.
3. Loop: `ok, err := f.TryLock()`. If ok, go to step 5. Otherwise sleep 100 ms (stop early on ctx done). After 1 s of waiting, read `info` once (ignore errors) and call `OnWait(holder or nil)`.
4. Deadline reached → `LockTimeout`, message `repository is locked by <command> (pid <pid> on <host>) since <started>`, or `repository is locked by another githive command` when the info file is missing or unreadable; hint `try again, or check that process`. Parent ctx cancelled → `Interrupted`.
5. Write the info file atomically: write `info.tmp-<pid>`, then rename it over `info`. If writing fails, unlock and return `IOFailed`.
6. `Release`: remove `info` (ignore not-exist), then `f.Unlock()`. Never delete `githive.lock` itself: deleting it while another process waits on it would let two processes lock two different files.

The CLI prints the wait message on stderr, in text mode only: `waiting for <command> (pid <pid>) on <project>…`.

## 6. Locking & concurrency

- The lock is per repository (`git rev-parse --git-common-dir`). Different repositories never wait for each other.
- Which steps take it is decided by each command's DD. DD-00 only provides `lock.With`. The rule, from BD-00 BR-09: fetch, worktree add/remove, branch delete, and `git config --local` writes.
- The OS releases a `flock`/`LockFileEx` lock when the process dies, so a crash leaves at most a stale `githive.lock.info`. That file is overwritten by the next holder, and it is only shown in a wait message.
- `git.Runner` is safe for concurrent use after the first call (`sync.Once` guards the version check).

## 7. Errors & exit codes

| Condition | Code | Message (text mode) | Exit |
| --- | --- | --- | --- |
| Unknown flag or command | `usage` | cobra's message, e.g. `unknown flag: --jsn` | 2 |
| Wrong number of arguments | `usage` | `accepts 0 arg(s), received 1` (cobra) | 2 |
| `GITHIVE_LOCK_TIMEOUT` invalid | `usage` | `GITHIVE_LOCK_TIMEOUT must be a whole number of seconds between 1 and 3600, got "<v>"` | 2 |
| git not on PATH | `git_missing` | `git not found on PATH` | 2 |
| git too old | `git_too_old` | `git 2.38 or newer is required, found <raw>` | 2 |
| git exits non-zero | `git_failed` | `git <cmd> failed: <first stderr line>` | 1 |
| Lock not acquired in time | `lock_timeout` | §5.3 step 4 | 5 |
| Interrupt (Ctrl+C, SIGTERM) | `interrupted` | `interrupted` | 130 |
| Writing the lock info file fails | `io_failed` | `cannot write <path>: <os error>` | 1 |
| Any other error (bug) | `internal` | the error text | 1 |

## 8. Platform specifics

- **Windows:** `gofrs/flock` uses `LockFileEx` on the whole file. That is why holder details live in the separate `githive.lock.info` and are never read from the locked file. `IsTerminal` works for Windows Terminal and conhost. `SIGTERM` does not exist, so only `os.Interrupt` is handled. `GCM_INTERACTIVE=never` stops Git Credential Manager from opening a login window for agents.
- **WSL:** `IsWSL` uses `WSL_DISTRO_NAME`, or `/proc/sys/kernel/osrelease` containing `microsoft`. `UnderMnt` only matters there.
- **macOS:** git from Xcode prints `git version 2.39.3 (Apple Git-146)`; the parser accepts it.
- **Linux:** none.

## 9. Test cases

| TC | Given / When / Then | Type | OS |
| --- | --- | --- | --- |
| TC-01 | `ExitCode` for every `Code` constant → value in §3.1; unknown code → 1 | unit | all |
| TC-02 | `apperr.From`: `*Error` unchanged; `context.Canceled` → `interrupted`; plain error → `internal` | unit | all |
| TC-03 | Success envelope with nil warnings → JSON has `"warnings": []`, no `error` key | unit | all |
| TC-04 | Failure envelope → `ok:false`, `error.code`, `error.exit_code`, no `data` key | unit | all |
| TC-05 | `ComputeMode` truth table: all 16 TTY × JSON × Yes combinations → `Interactive` only when both TTYs and neither flag | unit | all |
| TC-06 | `ComputeMode` colour: `NO_COLOR`, `--no-color`, `--json`, non-TTY each turn colour off | unit | all |
| TC-07 | `githive version --json` → one envelope with `command:"version"` and data `{version,commit,date}`; stderr empty; exit 0 | unit (Execute with buffers) | all |
| TC-08 | `githive --jsn` → exit 2; text mode prints `githive: unknown flag: --jsn` to stderr | unit | all |
| TC-09 | `githive --json --bogus` → exit 2; stdout is one envelope with `error.code:"usage"`; stderr empty | unit | all |
| TC-10 | `githive -y version` and `-v` accepted (short flags exist) | unit | all |
| TC-11 | Execute with an already-cancelled context and a command that returns `ctx.Err()` → exit 130, code `interrupted` | unit | all |
| TC-12 | `ParseVersion` table: plain, `.windows.1`, `(Apple Git-146)`, two-part `2.38`, garbage → error | unit | all |
| TC-13 | Runner with a fake version `2.37.1` (test hook) → first `Run` returns `git_too_old` | unit | all |
| TC-14 | Runner with `PATH` set to an empty temp dir → `git_missing` | integration | all |
| TC-15 | Non-interactive Runner: the env built for a call (exposed by a test hook) contains `LC_ALL=C`, `GIT_TERMINAL_PROMPT=0` and `GCM_INTERACTIVE=never`; an interactive Runner has `LC_ALL=C` but not the other two | unit | all |
| TC-16 | Runner on a failing command (`git rev-parse --verify nosuchref`) → `git_failed`, message starts `git rev-parse failed:` | integration | all |
| TC-17 | `redact` table: `https://user:tok@host/x` → `https://***@host/x`; URLs without credentials unchanged | unit | all |
| TC-18 | Runner with `Trace` set → one trace line per call, matching `^git -C .* \(\d+ ms\)$` | integration | all |
| TC-19 | `lock.Timeout`: unset → 30 s; `"5"` → 5 s; `"0"`, `"abc"`, `"4000"` → `usage` | unit | all |
| TC-20 | Two goroutines `Acquire` the same common dir: the second waits, `OnWait` is called once with the first holder's PID, and it acquires within 200 ms after the first `Release` | integration | all |
| TC-21 | Holder keeps the lock; second `Acquire` with Timeout 1 s → `lock_timeout` (exit 5) whose message contains the holder's command | integration | all |
| TC-22 | After `Release`, `githive.lock.info` is gone and `githive.lock` still exists; `Release` twice → no error | integration | all |
| TC-23 | Cross-process: test re-executes its own binary (helper-process pattern) to hold the lock; the parent gets `lock_timeout`; after the child exits the parent acquires | integration | all |
| TC-24 | Different common dirs lock independently (both acquire at once) | integration | all |
| TC-25 | `UnderMnt` with `IsWSL` forced true: `/mnt/c/x` → true, `/mnt/data/x` → false, `/home/x` → false | unit | all |
| TC-26 | `NoPrompter(false)`: `Input`, `Select` and `Confirm` → `input_required`; `NoPrompter(true)`: `Confirm` → true, `Input` and `Select` still → `input_required` (BR-02) | unit | all |
| TC-27 | `NativeAbs` table: relative path → absolute under the working dir; `a/../b` cleaned; on Windows the result uses `\`, elsewhere `/` | unit | all |
| TC-28 | `lock.With` whose `fn` blocks until ctx is cancelled: after cancel, `With` returns `interrupted`, `githive.lock.info` is gone, and a second `Acquire` succeeds at once (BR-12) | integration | all |

The test hooks (fake version, exposed env) are unexported package variables, set only in `_test.go` files.

## 10. Acceptance criteria

- [ ] Every error that reaches the user is an `*apperr.Error` with a code from §3.1, and its exit code follows the table (TC-01, TC-02, TC-11).
- [ ] `--json` prints exactly one envelope on stdout for success and failure, with `warnings` always an array (TC-03, TC-04, TC-07, TC-09).
- [ ] Global flags `--json`, `--yes/-y`, `--verbose/-v`, `--no-color` exist on every command (TC-08, TC-10).
- [ ] Interactive mode and colour follow BR-01 and §3.3 exactly (TC-05, TC-06, TC-26).
- [ ] git is run with an argument list, `LC_ALL=C`, and no credential prompts when non-interactive; version ≥ 2.38 is enforced (TC-12 … TC-16).
- [ ] git errors never show credentials (TC-17), and `--verbose` traces every call (TC-18).
- [ ] The repository lock serialises GitHive processes on one repo, never blocks other repos, times out with code 5, and reports who holds it (TC-19 … TC-24).
- [ ] Ctrl+C exits 130 and releases the lock (TC-11, TC-28).
- [ ] Paths shown or put in JSON are absolute and OS-native (TC-27).
- [ ] `internal/platform` has no import of `internal/ui` or `internal/apperr` (checked in review).
- [ ] CI is green on ubuntu, windows and macos; `go test -race` passes on Linux and macOS.

## 11. Open questions

- [x] Timeout for `fetch` → none; git's own network settings apply. Revisit if agents report hangs.
- [x] Package name for errors → `internal/apperr` (ADR-0004).
- No open question blocks implementation.

Coverage note: BR-13 (English only, CL-05) is a review convention and has no test case.

## Revision history

| Date | Change |
| --- | --- |
| 2026-11-13 | Approved (P04) |
| 2026-11-13 | Rev 2 (P04b): `internal/platform` moved into DD-00a and no longer imports `ui` (`RepoWarnings` removed, `NativeAbs` added); TC-27, TC-28 added; DD-03 split into 03a/03b |
| 2026-11-14 | Rev 3 (PO-01): "Test cases" column in the slices table, issues take their acceptance criteria from it (PO-01-D1, D2); `lock.With` returns errors through `apperr.From` (PO-01-D4); TC-26 also covers `Input` and `Select` (PO-01-D5) |
