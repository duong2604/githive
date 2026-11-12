---
id: ADR-0001
title: Record architecture decisions
status: Approved
req: []
depends_on: []
owner: Danny (duong2604)
last_updated: 2026-11-13
---

# ADR-0001 — Record architecture decisions

## Context

GitHive is built by AI agents and reviewed by an advisor and the product owner. Decisions made after the basic design ([`BD-00-overview.md`](../basic-design/BD-00-overview.md)) must be written down so that every agent session can find them and nobody re-argues them.

## Decision

Record every architecture decision taken after BD-00 as an ADR in `docs/adr/NNNN-short-title.md`, created from [`docs/templates/ADR-template.md`](../templates/ADR-template.md).

- An ADR is required to deviate from an approved design, to change the command surface, flags, exit codes, JSON schema or config keys, or to add a dependency outside the default list (CLAUDE.md §8).
- Status values: Draft, In review, Approved, Superseded, Rejected (P01b-D1). Nothing is built on an ADR until it is Approved.
- An ADR is proposed in its own pull request. When an ADR conflicts with a design document, the newest Approved ADR wins.

## Alternatives considered

- Editing the design documents directly: loses the reasoning and the history of why something changed.
- Recording decisions only in pull request descriptions: hard to find and not versioned with the code.

## Consequences

- Agents read `docs/adr/` together with the design documents at the start of a task.
- Small overhead per decision, accepted in exchange for traceability.

## Status

Approved (2026-11-13), as part of the bootstrap prompt P02.
