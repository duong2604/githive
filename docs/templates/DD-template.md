---
id: DD-xx
title:
status: Draft            # Draft | In review | Approved | Superseded
req: [REQ-xx, NFR-xx]
depends_on: []
owner:
last_updated: YYYY-MM-DD
---

# DD-xx — <feature name>

## 1. Scope
<!-- Which BD-xx this implements and which BR-xx rules it covers; list anything from the BD deferred elsewhere. -->

## 2. Packages & files
<!-- Go packages and files to create or change (see CLAUDE.md §5 layout), one line each with its responsibility. -->

## 3. Types & interfaces
<!-- Go type and function signatures this DD introduces or changes. -->

## 4. Config keys / JSON output schema
<!-- New or used config keys (global TOML and githive.* git config) and the exact --json schema. -->

## 5. Algorithm
<!-- Exact git commands in execution order, with arguments, and what is parsed from each output. -->

## 6. Locking & concurrency
<!-- Which steps hold githive.lock, for how long, and what happens when several agents run at once. -->

## 7. Errors & exit codes
<!-- Table: condition | message | exit code (0-5 per CLAUDE.md §4). -->

## 8. Platform specifics
<!-- Windows / WSL / Linux / macOS differences in paths, file locks, permissions, and how they are handled. -->

## 9. Test cases
<!-- Table: TC-xx | Given / When / Then | Type (unit / integration / e2e) | OS. -->

## 10. Acceptance criteria
<!-- Checklist; every line points to the TC-xx that proves it. -->

## 11. Open questions
<!-- Unresolved details as a checklist; implementation must not start while any blocks it. -->
