---
id: ADR-0004
title: Foundation contract additions (interrupt exit code, lock timeout, error package, dependencies)
status: Approved
req: [NFR-01, NFR-02, NFR-05]
depends_on: [ADR-0002, ADR-0003]
owner: Advisor
last_updated: 2026-11-13
---

# ADR-0004 — Foundation contract additions

## Context

[BD-00](../basic-design/BD-00-foundation.md) and [DD-00](../detail-design/DD-00-foundation.md) add items to the public contract that [ADR-0002](0002-core-technical-decisions.md) fixed: an exit code, an environment variable, a wider lock scope, a JSON envelope with stable error codes, a new package and new dependencies. CLAUDE.md §8 requires an ADR for each of these. The P03 report listed them.

## Decision

1. **Exit code 130** on interrupt (Ctrl+C, SIGTERM), with error code `interrupted`. This is the shell convention (128 + SIGINT). The ADR-0002 list 0–5 is otherwise unchanged.
2. **JSON envelope** `{ok, command, data, error, warnings}` for every `--json` output, with a stable `error.code` string per BD-00 §3.4. `githive version --json` moves to the envelope before v0.1.0.
3. **`GITHIVE_LOCK_TIMEOUT`** (whole seconds, 1–3600, default 30) overrides the lock timeout.
4. **Lock scope** includes `git config --local` writes, in addition to fetch, worktree add/remove and branch delete.
5. **Lock holder file** `githive.lock.info` next to `githive.lock`. The locked file itself is never read, because Windows `LockFileEx` blocks reads from other processes.
6. **New package** `internal/apperr` holds the error type and code catalogue shared by every package.
7. **New dependency** `golang.org/x/term` (terminal detection). `gofrs/flock` and `charmbracelet/huh` were already on the default list.
8. **Global flags** `--yes/-y`, `--verbose/-v`, `--no-color` join `--json` on every command.

## Alternatives considered

- Map interrupt to exit 1: agents could not tell a cancelled run from a failed one.
- Keep holder details inside `githive.lock`: unreadable on Windows while locked.
- Put error codes in `internal/ui`: services would then depend on the presentation package.
- Detect terminals with `os.File.Stat` mode bits: unreliable on Windows consoles; `golang.org/x/term` is the standard approach.

## Consequences

- CLAUDE.md §4 lists exit code 130, the global flags and the envelope; §5 lists `internal/apperr`, the lock info file and `golang.org/x/term`.
- The P02 `version --json` output changes shape (pre-v0.1.0, allowed).
- Later DDs pick error codes from the BD-00 catalogue and may add new ones; existing codes are never renamed or reused.

## Status

Approved by the advisor (2026-11-13) in P04, after Danny approved BD-00 and BD-03.
