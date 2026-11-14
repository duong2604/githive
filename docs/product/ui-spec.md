---
id: UI-SPEC
title: GitHive — CLI user interface specification
status: Approved
req: [REQ-01, REQ-02, REQ-03, REQ-04, REQ-05, REQ-06, REQ-07, REQ-08, NFR-01, NFR-04, NFR-05]
depends_on: [BD-00-overview, BD-00, BD-03]
owner: PO
last_updated: 2026-11-14
---

# GitHive — CLI user interface specification

GitHive has no window and no web page. Its interface is the terminal: bash or zsh on Linux, macOS and WSL; PowerShell, cmd and Git Bash on Windows; the terminal panel inside VS Code. This document describes everything a user sees and types there, from installing the tool to removing a finished task. Technical behaviour (flags, JSON, exit codes) is in BD-00 and BD-03; this document decides how it looks and feels.

## 1. Who uses the interface

| User | How they use GitHive | What matters most |
| --- | --- | --- |
| Danny (developer, several projects) | Types commands or uses the menu, several times a day | Few keystrokes, clear state of every task, no fear of losing work |
| AI agent (Claude Code) | Runs commands with `--json`, never interactive | Stable output, exit codes, no prompts, no hanging |
| A teammate trying the tool | Installs it, runs setup once | Install in one command, setup in under a minute, sensible defaults |

## 2. Design principles

- **P1 — Two modes, one tool.** A person in a terminal gets menus, colours and confirmations. A script or agent gets plain, stable output and never a question (BD-00 BR-01, BR-02).
- **P2 — Never lose work.** Anything that deletes shows exactly what will go and asks once (CL-03).
- **P3 — Show the next step.** Every result ends with what the user can do next: a `cd` command, the command to fix an error.
- **P4 — Fit any terminal.** Output fits 80 columns and works without colour or Unicode.
- **P5 — Same words everywhere.** One name per thing: project, task, base branch, worktree. Messages in English (CL-05).
- **P6 — This document is the reference for what users see.** Output examples in BD and DD documents are illustrative.

## 3. Interface surfaces

| ID | Surface | When the user sees it | Release |
| --- | --- | --- | --- |
| UI-01 | Install and first check | Once per machine | MVP (binary download, `go install`); v1 (winget, Scoop, Homebrew, install script) |
| UI-02 | First-run setup wizard (`githive init`) | First use, or to change settings | MVP |
| UI-03 | "Not set up yet" gate | Any command before setup | MVP |
| UI-04 | Main menu (`githive` with no arguments) | Daily use without remembering commands | MVP |
| UI-05 | Pickers: project, task, confirm | Inside the menu, or when a command lacks an argument in a terminal | MVP |
| UI-06 | Command results: `spawn`, `swarm`, `inspect`, `retire`, `project list`, `path`, `doctor` | Every command | MVP |
| UI-07 | `sync` screens | Updating a task from its base branch | v1 |
| UI-08 | Messages: errors, warnings, waiting | When something goes wrong or waits | MVP |
| UI-09 | Help text | `--help` on any command | MVP |
| UI-10 | Agent mode (`--json`) | Agents and scripts | MVP |
| UI-11 | Shell completion | Pressing Tab | v1 |

## 4. UI-01 Install and first check

MVP (release binaries and Go users):

```
$ go install github.com/duong2604/githive/cmd/githive@latest
$ githive version
githive 0.1.0 (commit 1a2b3c4, built 2026-12-01)
$ githive doctor
```

v1 (one command per OS):

```
Windows   winget install duong2604.githive      or   scoop install githive
macOS     brew install duong2604/tap/githive
Linux/WSL curl -fsSL https://.../install.sh | sh
```

After installing, the README tells the user to run `githive doctor`, then `githive init`.

## 5. UI-02 Setup wizard

Runs with `githive init`, or from the gate (UI-03). Every question shows a default; pressing Enter accepts it. The defaults follow BD-00-overview §4.

```
$ githive init

  GitHive setup  (step 1 of 4)

? Where are your projects?
  > D:\workspace                         (found 3 git projects)
    C:\Users\duong\workspace
    Type another path…

? Where should task worktrees go?   (step 2 of 4)
  > D:\workspace\worktrees               (same drive as your projects)
    Type another path…

? Branch prefix for new tasks?       (step 3 of 4)
  > feat/

? Base branch per project            (step 4 of 4)
  kirindo   dev        (found origin/dev)
  crm2      develop    (found origin/develop)
  legacy    —          not set — choose now? (y/N)

  Saved to C:\Users\duong\AppData\Roaming\githive\config.toml

  Next:  githive spawn kirindo@task-18     create your first task
         githive                           open the menu
```

Rules:

- On WSL the first suggestion is `~/workspace`, never `/mnt/...` (BD-00-overview §4).
- If the chosen projects folder has no git project, the wizard says so and asks again; it does not fail.
- Re-running `init` shows the current values as defaults and changes only what the user edits.
- Base branch detection only suggests a value; the user confirms each one.
- The MVP wizard does not ask for a sync strategy; the default is kept until sync arrives in v1 (UI-D8).

## 6. UI-03 "Not set up yet" gate

In a terminal:

```
$ githive swarm
GitHive is not set up on this machine yet.
? Run setup now? (Y/n)
```

Without a terminal, or with `--json`/`--yes`: error `not_configured`, exit 2, hint `githive init`.

Commands that work before setup, without the gate: `init`, `doctor`, `version`, `help`, `completion` (UI-D9).

## 7. UI-04 Main menu

`githive` with no arguments, in a terminal, opens the menu. Outside a terminal it prints help (unchanged from P01d-D3). The header shows where the user is.

```
$ githive

  GitHive 0.1.0   D:\workspace   3 projects · 5 tasks
  You are in: kirindo (main repo)

? What do you want to do?
  > New task                     create a worktree from the latest base branch
    My tasks                     list tasks and their state
    Open a task                  print the path to cd into
    Remove finished tasks        free disk space safely
    Projects and base branches
    Settings                     re-run setup
    Check environment            githive doctor
    Quit
```

Rules:

- Arrow keys or `j`/`k` to move, Enter to choose, Esc or `q` to go back or quit, `/` to filter long lists.
- When the user is inside a project, that project is pre-selected in every picker.
- After an action completes, the menu shows the same screen the command would show (UI-06), then returns to the menu. "Quit" or Esc on the main menu exits.
- Every menu action also prints the equivalent command in dim text, e.g. `→ githive spawn kirindo@task-18`, so users learn the commands.
- "Sync a task" appears in v1.
- In a command outside the menu, Esc or `q` in a picker cancels: GitHive prints `cancelled` and exits 130 (interrupted) (UI-D6).

## 8. UI-05 Pickers

Project picker shows the base branch so a project without one is obvious:

```
? Project
  > kirindo     base: dev       2 tasks
    crm2        base: develop   3 tasks
    legacy      base: NULL
```

`NULL` is the only word for a missing base branch, everywhere (UI-D2).

Choosing a project with no base branch asks for it first (same question as the wizard, step 4).

Task name input validates while typing and previews the result:

```
? Task name   task-18
  branch  feat/task-18
  folder  D:\workspace\worktrees\kirindo\task-18
```

An invalid name shows the reason and the suggestion under the field (BD-03 §3.3), e.g. `use lowercase — try task_18`. Enter is refused until the name is valid.

Task picker (for open, remove, sync) shows state so the right task is easy to spot:

```
? Task
  > task-18     feat/task-18     dirty    2h ago
    fix-otp     feat/fix-otp     merged   9d ago
```

## 9. UI-06 Command results

### spawn

```
$ githive spawn kirindo@task-18
✓ Fetched origin/dev (3f9c2a1, 12 minutes ago)
✓ Created branch feat/task-18
✓ Created worktree D:\workspace\worktrees\kirindo\task-18
✓ Copied .env
✓ Ran post-create: npm ci (41 s)

  Next:  cd D:\workspace\worktrees\kirindo\task-18
         claude
```

A failed hook shows `! post-create failed (exit 1)` with the last 5 lines of its output, keeps the worktree, and ends with the same "Next" block.

### swarm (task list)

```
$ githive swarm
kirindo   base origin/dev   2 tasks
  TASK       BRANCH          AHEAD  BEHIND  STATE    LAST COMMIT
  task-18    feat/task-18        3      12  dirty    2h ago
  task-19    feat/task-19        0       0  clean    5m ago
```

- STATE values: `clean`, `dirty` (uncommitted changes), `unpushed` (commits not on the remote), `merged` (remote branch gone or merged), `locked`, `external` (outside the worktrees folder, BD-03 BR-09).
- When several states apply, the first in this order is shown: `locked`, `dirty`, `unpushed`, `merged`, `clean`; `external` replaces all others. `inspect` lists every state that applies (UI-D3).
- BEHIND above 20 is shown in yellow with a footnote: `task-18 is 25 commits behind dev — consider githive sync` (v1 wording points to sync; MVP to `git rebase`).
- `--all` prints one block per project, separated by a blank line.
- Columns that do not fit narrow terminals are dropped from the right: LAST COMMIT first, then BRANCH.

### inspect

```
$ githive inspect kirindo@task-18
kirindo@task-18
  folder    D:\workspace\worktrees\kirindo\task-18
  branch    feat/task-18  (not pushed yet)
  base      origin/dev — 3 ahead, 12 behind
  changes   2 modified, 1 untracked
  created   2026-11-13 09:15 from 3f9c2a1
```

### retire (remove)

```
$ githive retire 'kirindo@*'
  TASK       STATE     SIZE     ACTION
  fix-otp    merged    1.2 GB   remove
  task-18    dirty     —        keep (2 uncommitted changes)
  task-19    unpushed  —        keep (3 commits not pushed)

? Remove 1 task and free 1.2 GB? (y/N)
✓ Removed kirindo@fix-otp
  Kept 2 tasks. Use --force to remove them anyway.
```

The default answer is No. Tasks kept for safety are always listed with the reason (CL-03).

Filters `--merged` and `--stale` arrive in v1 (UI-D1). Exit code: 0 when every selected task was removed; 4 (unsafe) when at least one selected task was kept by a safety check; tasks left out by a filter do not count. JSON lists `removed` and `kept`, each kept task with its reason (UI-D4).

### project list

```
$ githive project list
  PROJECT   BASE      TASKS  PATH
  kirindo   dev           2  D:\workspace\kirindo
  crm2      develop       3  D:\workspace\crm2
  legacy    NULL          0  D:\workspace\legacy
  → set a base branch: githive project config legacy --base <branch>
```

`NULL` is shown exactly as written in REQ-08.

### path

Prints only the path, so it can be used as `cd "$(githive path kirindo@task-18)"`. No colour, no extra lines.

### doctor

```
$ githive doctor
✓ git 2.46.0 (2.38 or newer required)
✓ config   C:\Users\duong\AppData\Roaming\githive\config.toml
✓ projects D:\workspace (3 found)
! core.longpaths is off — run: git config --global core.longpaths true
✗ repository kirindo is under /mnt/d — clone it inside WSL for speed
  2 problems found
```

## 10. UI-07 sync (v1)

```
$ githive sync
kirindo@task-18 is 12 commits behind origin/dev
  Checking for conflicts… none
? Update with: (use arrows)
  > rebase   (branch not pushed yet — recommended)
    merge
✓ Rebased feat/task-18 onto origin/dev (3 commits replayed)
```

When a conflict is predicted, the screen lists the files and offers "Cancel" (default), "Merge anyway" or "Keep conflicts to resolve by hand". The recommendation text follows BD-00-overview §9.

## 11. UI-08 Messages

- **Error:** `githive: <what went wrong>` then `hint: <exact command or action>`. Red `githive:` prefix in colour terminals.
- **Warning:** `! <message>` in yellow. Warnings never stop the command.
- **Success lines:** `✓` in green.
- **Waiting for the lock:** `waiting for githive spawn (pid 4312) on kirindo… (Ctrl+C to stop)`, updated once per second at most.
- **Long steps** (fetch, hooks) show a spinner with elapsed time in a terminal, and nothing extra in plain mode.
- Messages say what happened and what to do; they never blame the user and never show stack traces (unless `--verbose`).

## 12. Visual language

| Element | Colour terminal | No colour / no Unicode |
| --- | --- | --- |
| Success | green `✓` | `ok` |
| Warning | yellow `!` | `warning:` |
| Error | red `✗` / `githive:` | `error:` / `githive:` |
| Paths and commands to copy | cyan | plain |
| Secondary info (times, hints, equivalent commands) | dim | plain |
| Selected menu item | `>` and bold | `>` |

- Colour is off with `--no-color`, `NO_COLOR`, `--json`, or when output is not a terminal (BD-00).
- Unicode symbols fall back to the ASCII words above when the terminal cannot show them (Windows cmd with a legacy code page).
- Output that is not a terminal (pipe, file) uses the ASCII words (UI-D5).
- Known limitation: a standalone Git Bash window (mintty) may not be detected as a terminal; GitHive then runs in plain mode and prints a one-line tip: use Windows Terminal, the VS Code terminal, or `winpty githive` (UI-D10).
- Meaning never depends on colour alone: every coloured element also has a symbol or word.
- Tables fit 80 columns. Long paths are shortened in the middle (`D:\workspace\…\task-18`) in tables, never in `path` or JSON.
- Times are relative in tables (`5m ago`) and absolute in `inspect` and JSON.

## 13. UI-09 Help text

Every command's help follows one layout:

```
Create a task worktree from the latest base branch.

Usage:
  githive spawn <project@task> [flags]

Examples:
  githive spawn kirindo@task-18
  githive spawn task-18              (inside the kirindo repository)
  githive spawn kirindo@hotfix-1 --branch hotfix/LOGIN-42

Flags:
  --branch string   full branch name instead of <prefix><task>
  ...
Global flags:
  --json  --yes  --verbose  --no-color
```

Examples come first after Usage, because they are what users read.

## 14. UI-10 Agent mode

Agents and scripts use `--json` (and `--yes` where needed). Then: no colour, no spinner, no menu, no question, one JSON envelope on stdout, the exit code of ADR-0002/0004. The README has a short "For AI agents" section listing these rules and the error codes. Nothing in this document changes BD-00 §3.3.

## 15. UI-11 Shell completion (v1)

`githive completion bash|zsh|fish|powershell` prints a completion script. Tab completes commands, flags, project names and `project@task` targets. The `completion` command itself is available from the MVP (UI-D7).

## 16. Out of scope

- A full-screen dashboard (TUI that redraws the whole terminal): v2 at the earliest.
- Mouse input, desktop notifications, a GUI or web page.
- Languages other than English (CL-05).

## 17. Decisions by Danny

Answered 2026-11-13; recorded as clarifications CL-07 … CL-10 in `docs/requirements/README.md`.

- **CL-07 (NFR-05).** The main menu (UI-04) and pickers (UI-05) are in the MVP.
- **CL-08 (REQ-07).** Folder size is measured only in `retire`; `swarm` does not show size.
- **CL-09 (REQ-04).** After `spawn`, GitHive prints the `cd` command (and `claude`); opening a terminal tab or starting an agent automatically is v2.
- **CL-10 (REQ-07).** The removal confirmation defaults to No (y/N).

## 18. Traceability

| Requirement | Screens |
| --- | --- |
| REQ-01, REQ-02 (choose folders) | UI-02 steps 1–2 |
| REQ-03 (choose project) | UI-04, UI-05 project picker |
| REQ-04 (create and name a task) | UI-05 task name, UI-06 spawn |
| REQ-05 (list tasks) | UI-06 swarm, inspect |
| REQ-06 (update from base) | UI-07 sync (v1) |
| REQ-07 (remove several tasks) | UI-06 retire |
| REQ-08 (projects and base branch, NULL) | UI-06 project list, UI-02 step 4 |
| NFR-01 (Windows, WSL, Linux, macOS) | §12 fallbacks, UI-06 doctor |
| NFR-04 (settings survive restarts) | UI-02 "Saved to …" |
| NFR-05 (choose instead of typing) | UI-02, UI-04, UI-05 |

## 19. PO decisions (PO-02b)

Answered 2026-11-14 in `docs/prompts/PO-02b-ui-spec-decisions.md`.

- **UI-D1.** The MVP has no `--merged`/`--stale`; the `retire` example uses MVP syntax and both filters arrive in v1.
- **UI-D2.** One word, `NULL`, is used everywhere for a missing base branch (pickers, `project list`, JSON), as in REQ-08.
- **UI-D3.** STATE shows the first that applies in the order `locked`, `dirty`, `unpushed`, `merged`, `clean`; `external` replaces all others; `inspect` shows every state that applies.
- **UI-D4.** `retire` exits 0 when every selected task was removed and 4 (unsafe) when at least one selected task was kept by a safety check; tasks left out by a filter (v1) do not count; JSON lists `removed` and `kept`, each kept task with its reason.
- **UI-D5.** Output that is not a terminal (pipe, file) uses the ASCII words (`ok`, `warning:`, `error:`).
- **UI-D6.** Esc or `q` in a picker during a command (not in the menu) cancels: GitHive prints `cancelled` and exits 130 (interrupted).
- **UI-D7.** The cobra `completion` command is kept from the MVP; Tab completion of project and task names is v1.
- **UI-D8.** The MVP wizard has 4 steps and does not ask for a sync strategy (the default is kept until v1).
- **UI-D9.** Commands that run before setup, without the gate: `init`, `doctor`, `version`, `help`, `completion`.
- **UI-D10.** If a standalone Git Bash window (mintty) is not detected as a terminal, GitHive runs in plain mode and prints a one-line tip (Windows Terminal, the VS Code terminal, or `winpty githive`); a known limitation, to be confirmed after DD-00a.

## Revision history

| Date | Change |
| --- | --- |
| 2026-11-13 | In review (PO-02) |
| 2026-11-14 | Approved (PO-02b): UI-D1…UI-D10 |
