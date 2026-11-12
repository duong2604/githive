---
id: ADR-0003
title: Meaning of --all and project@*
status: Approved
req: [REQ-05, REQ-06, REQ-07]
depends_on: [ADR-0002]
owner: Advisor
last_updated: 2026-11-13
---

# ADR-0003 — Meaning of `--all` and `project@*`

## Context

The P01c consistency review found two readings of `--all`:

- [ADR-0002](0002-core-technical-decisions.md) and CLAUDE.md §4 say `project@*` **or** `-p project --all` selects every task of one project.
- [`BD-00-overview.md`](../basic-design/BD-00-overview.md) §7 shows `githive swarm [-p project | --all]`, where `--all` means every project.

The target grammar is a public interface (CLAUDE.md §4), so the meaning must be fixed before BD-03 and BD-05 are written.

## Decision

- `githive swarm --all` lists the worktrees of **every registered project**.
- `'project@*'` (quoted in shells) selects **every task of one project**. It is accepted by `swarm`, `sync` and `retire`.
- The form `-p project --all` is dropped.

This ADR clarifies ADR-0002; it does not supersede it. ADR-0002 stays Approved and unchanged, and is read together with this ADR.

## Alternatives considered

- Keep both `project@*` and `-p project --all` for "every task of one project": two spellings for the same selection, and `--all` would mean different things on different commands.
- Use `--all` only for "every task of one project": leaves no flag for listing every project, which REQ-05 and the multi-project positioning need.

## Consequences

- CLAUDE.md §4 drops `-p project --all` from the target grammar and documents `swarm --all`.
- BD-03 (target grammar) and BD-05 (`swarm`) define `--all` and `project@*` as above; DD-03 parses `project@*` as one target.
- Whether `-p/--project` remains as a filter on `swarm` (as written in BD-00 §7) is left to BD-05.

## Status

Approved by the advisor (2026-11-13) in P01d-D1.
