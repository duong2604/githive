---
id: ADR-0002
title: Core technical decisions
status: Approved
req: [NFR-01, NFR-02, NFR-03, NFR-04]
depends_on: [ADR-0001]
owner: Danny (duong2604)
last_updated: 2026-11-13
---

# ADR-0002 — Core technical decisions

## Context

[`BD-00-overview.md`](../basic-design/BD-00-overview.md) was written while the tool was called `wtx`. The project was renamed to GitHive and the command names changed. This ADR records the core technical choices as summarised in CLAUDE.md §4–5, so later ADRs can refer to them.

## Decision

- **Language and build:** Go, pinned in `go.mod`, built as a single static binary with GoReleaser for Windows, Linux and macOS on amd64 and arm64.
- **Git access:** shell out to the system `git` (minimum 2.38) with an argument slice, never through a shell; set `LC_ALL=C`; parse `--porcelain -z` output. No go-git.
- **State:** git is the source of truth for worktrees (`git worktree list --porcelain -z`); there is no state database.
- **Configuration, two layers:** global TOML at `%APPDATA%\githive\config.toml` on Windows or `$XDG_CONFIG_HOME/githive/config.toml` elsewhere; per-repo settings in `git config --local githive.*`. Files are written atomically. Precedence: flag > `GITHIVE_*` environment > repo git config > global config > defaults.
- **Locking:** `githive.lock` in `git rev-parse --git-common-dir`, held only around fetch, worktree add/remove and branch delete, with a 30 s timeout. Different repos never wait for each other.
- **Paths:** worktrees live at `<worktrees_root>/<project>/<task>`, never inside the repo; the default `worktrees_root` is `<projects_root>/worktrees`.
- **Target grammar:** `project@task`; `task` alone inside a registered repo or its worktrees; `project@*` or `-p project --all` for every task. Task names match `^[a-z0-9][a-z0-9._-]{0,49}$` and pass `git check-ref-format --branch`. No leading `@` (PowerShell splatting).
- **Commands** (BD-00 name → GitHive name): `wtx new` → `githive spawn` (aliases `new`, `create`); `wtx ls` → `githive swarm` (`ls`, `list`); `wtx rm` → `githive retire` (`rm`, `remove`); `project` has the alias `hive`. Packaging also installs `git-hive` so `git hive <cmd>` works.
- **Automation contract:** every command accepts `--json` and `--yes` and never waits for input without a TTY. Exit codes: 0 ok, 1 error, 2 bad arguments or missing config, 3 sync conflict (aborted cleanly), 4 blocked by a safety check, 5 lock timeout. JSON schemas and exit codes are a public API from v0.1.0.

## Alternatives considered

Covered in BD-00 section 3: Node/TypeScript and Rust instead of Go; go-git instead of the git CLI; a state file or SQLite instead of reading git; environment variables or a global project map instead of repo git config; no locking.

## Consequences

- Behaviour matches the user's own git (credentials, hooks, config) on every OS.
- Each command needs integration tests against real git (`internal/testutil`).
- Changing a command name, flag, exit code, JSON schema or config key requires a new ADR.

## Status

Approved (2026-11-13), as part of the bootstrap prompt P02. Summarises decisions already approved in BD-00 and CLAUDE.md §4–5.
