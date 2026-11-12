---
id: PROGRESS
title: Project progress
status: Approved
req: []
depends_on: []
owner: Danny (duong2604)
last_updated: 2026-11-13
---

# Progress

## Current phase

Phase 1 — the documentation baseline and the code skeleton are done (steps 1, 1b, 2a; 1c and 1d in review in PR #5). Next step: 2b (BD-00-foundation, BD-03 → DD-00, DD-03).

## Scope

From `BD-00-overview.md` section 10 (P01c-D4):

| Phase | Scope |
| --- | --- |
| MVP | `init`, `config`, `project add/list/config`, `spawn` (with `copy_files` and `postCreate`), `swarm`, `inspect`, `retire <targets>` with CL-03 safety checks, `path`, minimal `doctor`, `--json` and exit codes, per-repo locking, two-layer config |
| v1 | `sync` with conflict prediction, `retire --merged/--stale/-i`, squash-merge detection, `shell-init`, release with GoReleaser, Scoop/Homebrew and the `git-hive` alias; Windows archives as `.zip` via GoReleaser `format_overrides` (P01d-D7) |
| v2 | Forge API (MR status, create MR), `--run claude`, busy-worktree marking, per-worktree ports, shared `.githive.toml`, `agent-rules`, TUI dashboard |

## Plan

| Step | Owner | Status | Output | Date |
| --- | --- | --- | --- | --- |
| 1 Normalize documentation structure | Agent | Done | [PR #1](https://github.com/duong2604/githive/pull/1) | 2026-11-13 |
| 1b Apply P01 review decisions (P01b) | Agent | Done | [PR #3](https://github.com/duong2604/githive/pull/3) (merged) | 2026-11-13 |
| 1c Refresh BD-00 overview (P01c) | Agent | In review | [PR #5](https://github.com/duong2604/githive/pull/5) | 2026-11-13 |
| 1d Rebase on bootstrap, apply P01c review (P01d) | Agent | In review | [PR #5](https://github.com/duong2604/githive/pull/5) | 2026-11-13 |
| 2a Bootstrap project (P02, runs P00 part B steps 0–2) | Agent | Done | [PR #4](https://github.com/duong2604/githive/pull/4) (merged) | 2026-11-13 |
| 2b BD-00, BD-03 → DD-00, DD-03 approved | Advisor + Danny | Not started | | |
| 3 Issues + implementation for DD-00, DD-03 | Agent | Not started | | |

Prompts for each step: [`docs/prompts/README.md`](prompts/README.md) (one prompt per session).

## Review findings (2026-11-13)

| # | Finding | Status | Resolution |
| --- | --- | --- | --- |
| 1 | Missing architecture doc in the repo | Fixed | `docs/basic-design/BD-00-overview.md` added (Markdown export) |
| 2 | Binary documents (.docx, .xlsx) | Fixed | Converted to Markdown and removed; `docs-check` rejects non-.md files |
| 3 | No requirement IDs | Fixed | REQ-01..08, NFR-01..05 and a Traceability table in `docs/requirements/README.md` |
| 4 | Folder names with spaces, capitals and a typo (`Promtps`) | Fixed | Renamed to kebab-case ASCII |
| 5 | README.md was UTF-16 with CRLF | Fixed | Rewritten as UTF-8, LF; `.gitattributes` enforces LF |
| 6 | Documentation commits used `feat:` | Fixed | Commit type rule added to CLAUDE.md §7 (existing history unchanged) |
| 7 | Kickoff prompt created MVP issues before any DD existed | Fixed | Kickoff limited to steps 0–2; issues come from approved DDs (P03) |
| 8 | No shared templates or automated docs check | Fixed | `docs/templates/` and `go run ./tools/docs-check`; CI `docs-check` job added in P02 |
| 9 | Document status values inconsistent (Approved vs Proposed/Accepted) | Fixed | P01b-D1: one status set; Rejected for ADRs only |
| 10 | Date format deviated from the template (dd-mm-yyyy) | Fixed | P01b-D5: ISO `YYYY-MM-DD` everywhere; docs-check rule (e) |
| 11 | `go.work` at the repo root is temporary | Fixed | P01b-D2: removed in P02 when the root `go.mod` was created |
| 12 | No fixed report format in CLAUDE.md | Fixed | P01b-D6: CLAUDE.md §13 "Reporting & escalation" |
| 13 | Prompt files not identifiable per task | Fixed | P01b-D7: `P<NN>[x]-<slug>.md`, prompt front matter, `docs/prompts/README.md` index; docs-check rule (f) |

## Log

| Date | Work | PR | OS verified |
| --- | --- | --- | --- |
| 2026-11-13 | Normalized docs structure: kebab-case folders, requirements and BD index in Markdown, templates, docs-check tool, README encoding, CLAUDE.md and kickoff prompt updated | [#1](https://github.com/duong2604/githive/pull/1) | Windows 11 (Git Bash, Go 1.26.5) |
| 2026-11-13 | PR #1 merged into main (merge commit ee1bf78); branch deleted. Step 1 marked Done | — | Windows 11 |
| 2026-11-13 | P01b: applied P01b-D1…D7 and CL-01…04; docs-check rules (e) and (f); prompts split into P00–P03 with an index; PR #2 folded in and closed | [#3](https://github.com/duong2604/githive/pull/3) | Windows 11 (Git Bash, Go 1.26.5) |
| 2026-11-13 | P02: Go module, `githive version` + `--json`, `internal/testutil`, CI (3 OS + docs-check + goreleaser check), GoReleaser config, MIT license, README/CONTRIBUTING/CHANGELOG, ADR-0001/0002; docs-check merged into the root module | [#4](https://github.com/duong2604/githive/pull/4) | Windows 11 locally; CI green on ubuntu, windows, macos |
| 2026-11-13 | P01c: BD-00 overview revision 2 committed verbatim; wtx naming map dropped; `extra_projects` in CLAUDE.md §5; Scope section; consistency pass (mismatches reported, not changed) | [#5](https://github.com/duong2604/githive/pull/5) | Windows 11 (Git Bash, Go 1.26.5) |
| 2026-11-13 | P01d: PR #5 rebased on main after PR #4; ADR-0003 (`--all`, `project@*`); CLAUDE.md §4 (version, shell-init, retire --force, bare `githive` prints help) and §5 key naming; BD-00 overview vs foundation note; v1 zip note | [#5](https://github.com/duong2604/githive/pull/5) | Windows 11 (Git Bash, Go 1.26.5) |
