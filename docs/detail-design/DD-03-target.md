---
id: DD-03
title: Target grammar — parser and resolver
status: Approved
req: [REQ-03, REQ-04, REQ-05, REQ-07]
depends_on: [BD-03, DD-00, ADR-0003]
owner: Advisor
last_updated: 2026-11-13
---

# DD-03 — Target grammar

## 1. Scope

Implements [BD-03](../basic-design/BD-03-target.md), rules BR-01 … BR-11. It provides a pure parser and a resolver that later commands (spawn, swarm, inspect, sync, retire, path) call.

The resolver does not know where projects and tasks come from. It depends on two small interfaces (§3.3), implemented later by the project registry (DD-02) and the worktree service (DD-04/DD-05). Until those exist, tests use fakes.

### Implementation slices (one issue each)

| Slice | Contents | Depends on | Test cases |
| --- | --- | --- | --- |
| DD-03a | `parse.go`, `resolve.go`, `spec.go`, `distance.go`; the resolver tested with fake sources, locator and prompter | DD-00a | TC-01 … TC-13, TC-19 … TC-22 |
| DD-03b | `locator.go` (`gitLocator` on `git.Runner`) and the integration tests with real repositories | DD-00b, DD-03a | TC-14 … TC-18 |

DD-03a runs in parallel with DD-00b and DD-00c.

Deferred elsewhere:
- Checking `<prefix><task>` with `git check-ref-format --branch`: DD-04 (`spawn`), because the prefix comes from project config.
- `project add --name`: DD-02.
- Listing external worktrees (BR-09): DD-05.

## 2. Packages & files

| Path | Responsibility |
| --- | --- |
| `internal/target/parse.go` | `Parse`, `ValidateTask`, `ValidateProject`, `SuggestTask` — pure functions, no I/O |
| `internal/target/resolve.go` | `Resolver`, `Spec`, `Resolved`; current-directory lookup; closest-name hints |
| `internal/target/spec.go` | One `Spec` value per command (§3.4) |
| `internal/target/distance.go` | Levenshtein distance for "did you mean" (about 25 lines, no dependency) |
| `internal/target/locator.go` | `gitLocator`, the default `Locator` on DD-00's `git.Runner` (slice DD-03b) |

## 3. Types & interfaces

### 3.1 Parsed target

```go
package target

// Target is what the user typed, before lookup.
type Target struct {
    Raw     string // as given, for messages
    Project string // "" → current project (bare task form)
    Task    string // "" only when All is true
    All     bool   // project@*
}

// Parse accepts "project@task", "task" or "project@*".
func Parse(s string) (Target, error) // error is *apperr.Error with code invalid_target

func ValidateTask(name string) error    // BD-03 §3.3; invalid_target with a hint from SuggestTask
func ValidateProject(name string) error // ^[A-Za-z0-9._-]{1,64}$
func SuggestTask(name string) string    // "" when no valid suggestion exists
```

### 3.2 Resolution result

```go
type Project struct {
    Name      string // as registered (display form)
    Path      string // main repo directory, absolute
    CommonDir string // absolute `git rev-parse --git-common-dir`
    TasksDir  string // <worktrees_root>/<Name>, absolute
}

type Resolved struct {
    Project Project
    Task    string // "" when the spec selects the whole project (list commands)
    Path    string // <TasksDir>/<Task>; "" when Task == ""
    Exists  bool   // the task worktree exists
}

// JSON shape used by every command that reports a task (BD-03 §3.4):
type Ref struct {
    Project string `json:"project"`
    Task    string `json:"task"`
    Target  string `json:"target"` // "project@task"
}
func (r Resolved) Ref() Ref
```

### 3.3 Interfaces the resolver needs

```go
// ProjectSource lists registered projects (DD-02 implements it).
type ProjectSource interface {
    Projects(ctx context.Context) ([]Project, error)
}

// TaskSource lists the task names of a project: directories under TasksDir
// that are worktrees of the project (DD-04/DD-05 implement it).
type TaskSource interface {
    Tasks(ctx context.Context, p Project) ([]string, error)
}

// Locator finds which project/task the current directory belongs to.
type Locator interface {
    // Locate returns the absolute git common dir and the worktree top level
    // for dir, or ok=false when dir is not inside a git repository.
    Locate(ctx context.Context, dir string) (commonDir, topLevel string, ok bool, err error)
}

type Resolver struct {
    Projects ProjectSource
    Tasks    TaskSource
    Locator  Locator
    Cwd      string
    Prompt   ui.Prompter // ui.NoPrompter when not interactive
}

// Resolve returns the results, warnings (repo_on_mnt) and an *apperr.Error.
func (r *Resolver) Resolve(ctx context.Context, args []string, spec Spec) ([]Resolved, []ui.Warning, error)
```

`gitLocator` is the default `Locator`, built on DD-00's `git.Runner` (§5.3).

### 3.4 Command specs

```go
type Implicit int

const (
    ImplicitNone    Implicit = iota // a target is required
    ImplicitTask                    // no target → current task
    ImplicitProject                 // no target → current project, Task = ""
)

type Spec struct {
    Implicit      Implicit
    AllowTask     bool // project@task and bare task forms
    AllowWildcard bool // project@*
    AllowMultiple bool
    MustExist     bool // false for spawn
}

var (
    SpecSpawn   = Spec{Implicit: ImplicitNone, AllowTask: true, MustExist: false}
    SpecInspect = Spec{Implicit: ImplicitTask, AllowTask: true, MustExist: true}
    SpecPath    = Spec{Implicit: ImplicitTask, AllowTask: true, MustExist: true}
    SpecSync    = Spec{Implicit: ImplicitTask, AllowTask: true, AllowWildcard: true, AllowMultiple: true, MustExist: true}
    SpecRetire  = Spec{Implicit: ImplicitNone, AllowTask: true, AllowWildcard: true, AllowMultiple: true, MustExist: true}
    SpecSwarm   = Spec{Implicit: ImplicitProject, AllowWildcard: true} // AllowTask false: swarm takes 'project@*' or --all
)
```

This table is BD-03 §3.2 in code. `swarm --all` is handled by the swarm command itself (DD-05), not by the resolver.

## 4. Config keys / JSON output schema

- **No config keys.** `worktrees_root` and the project list come through `ProjectSource`.
- **JSON:** `Ref` (§3.2) is embedded by commands as `"project"`, `"task"`, `"target"`.
- **Error hints** in JSON use the same `hint` field as text mode (DD-00 envelope).

## 5. Algorithm

### 5.1 `Parse(s)`

1. `s == ""` → `invalid_target` "empty target".
2. `strings.HasPrefix(s, "@")` → `invalid_target` `target must not start with "@"` (BR-03).
3. `n := strings.Count(s, "@")`. If `n > 1` → `invalid_target` `target has more than one "@"`.
4. If `n == 0`: `ValidateTask(s)`; return `{Raw: s, Task: s}`.
5. Split at `@` into `p`, `t`. Empty `p` or empty `t` → `invalid_target` `target "<s>" needs both a project and a task, like kirindo@task-18`.
6. `ValidateProject(p)`.
7. If `t == "*"` → return `{Project: p, All: true}`.
8. `ValidateTask(t)`; return `{Project: p, Task: t}`.

### 5.2 `ValidateTask` and `SuggestTask`

Rules, checked in this order (the first failure is reported):
1. Matches `^[a-z0-9][a-z0-9._-]{0,49}$`.
2. Does not contain `..`.
3. Does not end with `.` or `.lock`.

Error message: `invalid task name "<name>": use lowercase letters, digits, ".", "_" or "-" (max 50, must start with a letter or digit)`. Hint: `githive <cmd> <project>@<suggestion>` when `SuggestTask` returns a value. The CLI fills in `<cmd>`; the parser returns only the suggestion.

`SuggestTask(name)`:
1. Lowercase.
2. Replace each run of whitespace with `-`.
3. Delete every character outside `[a-z0-9._-]`.
4. Collapse repeated `-`.
5. Replace `..` with `.`; strip a trailing `.lock` or `.`.
6. Trim leading characters until the first `[a-z0-9]`.
7. Cut to 50 characters.
8. Return the result only if it passes `ValidateTask`; otherwise `""`.

Examples: `Task_18` → `task_18`; `Fix Login` → `fix-login`; `API--M018` → `api-m018`; `@@@` → `""`.

### 5.3 Current directory (`gitLocator.Locate`)

1. `git -C <dir> rev-parse --path-format=absolute --git-common-dir --show-toplevel` (needs git 2.31+, covered by the 2.38 minimum).
2. Failure with "not a git repository" in stderr → `ok=false`, no error. Any other failure → the `git_failed` error.
3. Output: line 1 = common dir, line 2 = top level. Clean both with `filepath.Clean`.

### 5.4 `Resolve(ctx, args, spec)`

1. **No args:**
   - `ImplicitNone` → `input_required` `a target is required, like kirindo@task-18`; in interactive mode, ask with `Prompt.Input` instead.
   - `ImplicitTask` → `current(ctx)` must return a project and a task. When it finds a project but no task (the user is in the main repo) → `input_required` (`not inside a task worktree; pass a target like <project>@<task>`), or in interactive mode `Prompt.Select` from that project's tasks. When it finds nothing at all → `input_required`, or in interactive mode `Prompt.Select` from the registered projects, then from that project's tasks.
   - `ImplicitProject` → `current(ctx)` must return a project, else `input_required` `not inside a project; use 'project@*' or --all`, or in interactive mode `Prompt.Select` from the registered projects. Returns one `Resolved` with `Task: ""`.
2. **Several args** when `!AllowMultiple` → `usage` `this command takes one target`.
3. For each arg: `Parse`. If `All` and `!AllowWildcard` → `invalid_target` `<cmd> does not accept project@*`. If `!All` and `!AllowTask` → `invalid_target` `<cmd> takes 'project@*' or --all, not a single task`.
4. **Project lookup:**
   - `Project == ""` → use the project from `current(ctx)`; if none → `input_required` `"<task>" has no project here; use <project>@<task>`, or in interactive mode `Prompt.Select` from the registered projects (BD-03 §5).
   - Otherwise find a registered project with `strings.EqualFold(name, p.Name)`. Not found → `project_not_found` `project "<p>" not found`, hint `did you mean "<a>", "<b>"? Registered projects: githive project list` (up to 3 names with distance ≤ 3, case-insensitive; omit "did you mean" when none).
5. **Task lookup:**
   - `All` → `Tasks.Tasks(p)`; one `Resolved` per task, sorted by name. An empty list is valid.
   - Otherwise `Exists = task ∈ Tasks.Tasks(p)`. If `MustExist && !Exists` → `task_not_found` `task "<t>" not found in project "<p>"`, hint with up to 3 close names.
6. `Path = filepath.Join(p.TasksDir, task)`.
7. **Duplicates:** drop a `Resolved` whose `(Project.Name, Task)` was already seen; keep first-seen order (BR-10).
8. **Warnings:** for every distinct project whose `platform.UnderMnt(p.Path)` is true, append `ui.Warning{Code: "repo_on_mnt", Message: "repository <p.Path> is under /mnt; git is slow there, clone it inside WSL instead"}`. Resolve returns them as its second result (§3.3). `internal/platform` stays free of `internal/ui` (DD-00 §3.6).

### 5.5 `current(ctx)`

1. `Locator.Locate(ctx, Cwd)`. If `!ok` → no project.
2. Find the registered project whose `CommonDir` is the same directory as the located common dir, compared with `os.SameFile` on `os.Stat` of both. This handles case, symlinks, `\` vs `/` and 8.3 short names on Windows.
3. Task: if `topLevel` is a direct child of `p.TasksDir` (same `os.SameFile` check on `filepath.Dir(topLevel)`), the task is `filepath.Base(topLevel)`. Otherwise there is no current task (main repo or external worktree).

## 6. Locking & concurrency

None. Resolution only reads: git `rev-parse`, the project list and directory listings. It never takes `githive.lock`. A task created or removed by another agent between resolution and use is handled by the command, which re-checks inside its own lock (DD-04, DD-07).

## 7. Errors & exit codes

| Condition | Code | Message | Exit |
| --- | --- | --- | --- |
| Empty target, leading `@`, more than one `@`, empty side | `invalid_target` | §5.1 | 2 |
| Bad task or project name | `invalid_target` | §5.2, hint with suggestion | 2 |
| `project@*` where not allowed | `invalid_target` | `<cmd> does not accept project@*` | 2 |
| Several targets where one is allowed | `usage` | `this command takes one target` | 2 |
| Target needed but none given (non-interactive) | `input_required` | §5.4 step 1 | 2 |
| Bare task outside any project | `input_required` | §5.4 step 4 | 2 |
| Unknown project | `project_not_found` | §5.4 step 4, hint with close names | 2 |
| Unknown task (when it must exist) | `task_not_found` | §5.4 step 5, hint with close names | 2 |
| `rev-parse` fails for a reason other than "not a git repository" | `git_failed` | from DD-00 runner | 1 |

## 8. Platform specifics

- **Windows:** `os.SameFile` compares volume serial and file index, so `D:\Workspace\kirindo\.git` and `d:/workspace/kirindo/.git` match. Project names are compared with `strings.EqualFold`.
- **WSL:** `/mnt/d/workspace/kirindo` (Windows-side clone) and `~/workspace/kirindo` are different repositories and never match. Resolution adds a `repo_on_mnt` warning for the first (§5.4 step 8).
- **macOS:** `os.SameFile` avoids comparing Unicode forms of path names. Project names are ASCII only (§3.1), so they need no normalisation.
- **Linux:** none.

## 9. Test cases

| TC | Given / When / Then | Type | OS |
| --- | --- | --- | --- |
| TC-01 | `Parse` valid table: `kirindo@task-18`, `task-18`, `kirindo@*`, `Kirindo@v2.1-hotfix`, `crm2@api_m018` → fields as expected | unit | all |
| TC-02 | `Parse` invalid table: `""`, `@task`, `kirindo@`, `a@b@c`, `kirindo@Task_18`, `kir indo@x`, `k/x@t` → `invalid_target` | unit | all |
| TC-03 | `ValidateTask` boundaries: 50 chars ok, 51 fails; `a..b`, `x.lock`, `x.`, `-x`, `_x` fail; `a`, `0`, `v2.1` ok | unit | all |
| TC-04 | `SuggestTask` table: `Task_18`→`task_18`, `Fix Login`→`fix-login`, `API--M018`→`api-m018`, `x.lock`→`x`, `@@@`→`""`, 60 chars → 50 | unit | all |
| TC-05 | Resolve `kirindo@task-18` with fake sources where the task exists → one result, `Exists` true, `Path` = `TasksDir/task-18` | unit | all |
| TC-06 | Resolve `KIRINDO@task-18` → matches project `kirindo`; result shows the registered name | unit | all |
| TC-07 | Resolve `kirndo@task-18` → `project_not_found`, hint contains `"kirindo"` | unit | all |
| TC-08 | Resolve `kirindo@task-81` with SpecInspect → `task_not_found`, hint contains `task-18`; with SpecSpawn → ok, `Exists` false | unit | all |
| TC-09 | SpecSpawn with `kirindo@*` → `invalid_target`; SpecRetire with `kirindo@*` and 3 tasks → 3 results sorted | unit | all |
| TC-10 | SpecRetire with `kirindo@*` and 0 tasks → empty result, no error | unit | all |
| TC-11 | SpecRetire with no args → `input_required` even when standing in a task worktree | unit | all |
| TC-12 | SpecInspect with two args → `usage` | unit | all |
| TC-13 | SpecSync with `kirindo@a crm2@b kirindo@a` → 2 results in that order | unit | all |
| TC-14 | Real repo (testutil fixture) registered as `kirindo` with TasksDir; a worktree added with `git worktree add <TasksDir>/task-18 -b feat/task-18 origin/dev`: Cwd = worktree → SpecInspect with no args resolves `kirindo@task-18` | integration | all |
| TC-15 | Same fixture, Cwd = main repo: bare `task-18` resolves; no-arg SpecInspect → `input_required` | integration | all |
| TC-16 | Cwd = a temp dir outside any repo: SpecSwarm with no args → `input_required` naming `'project@*'` and `--all` | integration | all |
| TC-17 | Cwd = a subdirectory deep inside the worktree → still resolves the task | integration | all |
| TC-18 | Windows: project `CommonDir` registered with different case and `/` separators still matches (SameFile) | integration | windows |
| TC-19 | Interactive resolver (fake Prompter) with no args in the main repo and SpecInspect → `Select` is offered the project's tasks; the chosen task is resolved | unit | all |
| TC-20 | `UnderMnt` hook forced to true → Resolve returns one `repo_on_mnt` warning per distinct project | unit | all |
| TC-21 | Interactive resolver (fake Prompter), Cwd outside any project, bare `task-18` → `Select` is offered the registered projects; the chosen project is used | unit | all |
| TC-22 | SpecSwarm with `kirindo@task-18` → `invalid_target` naming `'project@*'` and `--all`; with `kirindo@*` → ok | unit | all |

## 10. Acceptance criteria

- [ ] All three target forms parse, and every invalid form is rejected with `invalid_target` (TC-01, TC-02).
- [ ] Task name rules and suggestions follow BD-03 §3.3 and CL-06 (TC-03, TC-04).
- [ ] Project lookup is case-insensitive with "did you mean" hints (TC-06, TC-07).
- [ ] Each command's allowed forms match BD-03 §3.2 (TC-08 … TC-13, TC-22).
- [ ] In interactive mode a missing project or task is picked from a list, as BD-03 §5 requires (TC-19, TC-21).
- [ ] The current project and task are found from the current directory, including subdirectories and Windows path variants (TC-14 … TC-18).
- [ ] WSL `/mnt` warnings are attached once per project (TC-20).
- [ ] No I/O in `parse.go`; slice DD-03a passes `go test ./internal/target/...` without git (TC-01 … TC-13, TC-19 … TC-22).

## 11. Open questions

- None block implementation. `project add --name` (DD-02) and external worktrees (DD-05) only add data behind the interfaces in §3.3.

Coverage deferred to later DDs: BR-05 "two projects may not differ only in case" is enforced at registration (DD-02); BR-09 external worktrees are listed by `swarm` (DD-05); BR-11 quoting is a documentation rule with no test case.

## Revision history

| Date | Change |
| --- | --- |
| 2026-11-13 | Approved (P04) |
| 2026-11-13 | Rev 2 (P04b): split into DD-03a/DD-03b; `Spec.AllowTask` so `swarm` rejects single tasks; interactive project picking per BD-03 §5; `repo_on_mnt` warning built here instead of in `platform`; TC-21, TC-22 added |
