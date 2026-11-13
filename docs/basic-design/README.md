---
id: BD-INDEX
title: Basic design index
status: Approved
req: []
depends_on: []
owner: Danny (duong2604)
last_updated: 2026-11-13
---

# Basic design index

The BD-00 row below is the feature document "Shared foundation", written as [`BD-00-foundation.md`](BD-00-foundation.md); it covers NFR-01, NFR-02 and NFR-05. [`BD-00-overview.md`](BD-00-overview.md) is the overall architecture and therefore spans every REQ (P01d-D6).

| BD ID | Feature | Commands | DD ID | BD status | DD status | Depends on | REQ |
| --- | --- | --- | --- | --- | --- | --- | --- |
| BD-00 | [Shared foundation](BD-00-foundation.md): git runner, locking, `--json` output and exit codes, OS-specific paths | (shared) | DD-00 | In review | Not started | — | NFR-01, NFR-02, NFR-05 |
| BD-01 | Initialisation and global configuration | `init`, `config` | DD-01 | Not started | Not started | BD-00 | REQ-01, REQ-02, NFR-04, NFR-05 |
| BD-02 | Projects and base branch | `project add/list/config` | DD-02 | Not started | Not started | BD-00, BD-03 | REQ-03, REQ-08, NFR-03 |
| BD-03 | [`project@task` target grammar](BD-03-target.md) | (shared) | DD-03 | In review | Not started | — | REQ-03, REQ-04, REQ-05, REQ-07 |
| BD-04 | Create a worktree | `spawn` | DD-04 | Not started | Not started | BD-00, BD-03 | REQ-04, NFR-02 |
| BD-05 | Show status | `swarm`, `inspect` | DD-05 | Not started | Not started | BD-00, BD-03 | REQ-05 |
| BD-06 | Sync with base | `sync` | DD-06 | Not started | Not started | BD-00, BD-03 | REQ-06 |
| BD-07 | Safe cleanup | `retire` | DD-07 | Not started | Not started | BD-00, BD-03 | REQ-07 |
| BD-08 | Navigation and shell | `path`, `shell-init` | DD-08 | Not started | Not started | BD-00, BD-03 | — (design-derived) |
| BD-09 | Environment diagnostics | `doctor` | DD-09 | Not started | Not started | BD-00 | NFR-01 |

BD-00 and BD-03 are settled first: every other feature depends on them.
BD-08 has no source requirement; it exists so an agent can move into its worktree.
REQ mapping comes from the Traceability table in [`docs/requirements/README.md`](../requirements/README.md).

## Conventions

- **Basic design (BD) answers "what".** Goal, users, use cases, user-facing input/output, main and error flows, business rules, per-OS constraints, out of scope.
- **Detail design (DD) answers "how".** Packages and functions, structs, config and JSON schemas, the exact git commands in order, lock scope, error and exit-code table, test cases.
- Each feature has one BD/DD pair with the same number (BD-04 ↔ DD-04). Issues and PRs name the DD they implement.
- Traceability chain: REQ → BD-xx → DD-xx → issue → PR.
- **Code is written only from a DD whose status is Approved.**
- Status flow in this table: Not started → Draft → In review → Approved (→ Superseded). Document status values: Draft | In review | Approved | Superseded, plus Rejected for ADRs only (P01b-D1).
- File names: `BD-xx-<slug>.md` in `docs/basic-design/`, `DD-xx-<slug>.md` in `docs/detail-design/`. Each file starts from [`docs/templates/`](../templates/) and opens with the YAML front matter block (id, title, status, req, depends_on, owner, last_updated).
- Dates use ISO 8601 `YYYY-MM-DD` everywhere (P01b-D5).

## Naming

[`BD-00-overview.md`](BD-00-overview.md) uses the official GitHive names since revision 2 (P01c); no name mapping is needed.

The progress tracker lives on claude.ai; this repo keeps only the Markdown copy.
