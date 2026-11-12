---
id: REQ-FORMAT
title: Requirement ID format and traceability rules
status: Approved
req: []
depends_on: []
owner: Danny (duong2604)
last_updated: 13-11-2026
---

# Requirement ID format

Requirements live in [`docs/requirements/README.md`](../requirements/README.md). The original text is kept verbatim in the language it was written in; IDs are added as bold tags in front of the paragraph they come from, e.g. `**[REQ-04]**`.

## IDs

- `REQ-xx`: functional requirement, something the user can do. Two digits, starting at `REQ-01`.
- `NFR-xx`: non-functional requirement, a quality the whole tool must have (platforms, persistence, concurrency, UX). Two digits, starting at `NFR-01`.
- IDs are never reused or renumbered. A dropped requirement keeps its ID and gets status `Dropped` in the Traceability table.
- One paragraph may carry several tags; one ID may be tagged on several paragraphs.
- A new requirement gets the next free number and is appended to the original text in its own section with the date (`dd-mm-yyyy`) and the author.

## Requirement IDs table

Every ID has one row: `ID | Summary | Source | BD`. The summary is a one-line English paraphrase; the source is the paragraph number in the original text.

## Traceability table

Columns: `REQ | BD | DD | Issue | Status`.

- **BD / DD:** set when the BD is created; the DD number always matches the BD number.
- **Issue:** filled when the issue is created from an Approved DD (`#N`, several separated by commas).
- **Status:** `Not started` → `In design` (BD or DD in Draft/In review) → `Ready` (DD Approved) → `In progress` (issue open) → `Done` (all issues closed), or `Dropped`.
- Update the table in the same PR that changes the BD, DD or issue it reflects, and mirror the change in [`docs/basic-design/README.md`](../basic-design/README.md).
- A BD that comes from the design rather than a requirement is listed in a note under the table, with the reason it exists.
