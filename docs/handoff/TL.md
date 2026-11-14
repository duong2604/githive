# Handoff — TL

| Field | Value |
| --- | --- |
| Role / agent ID | Technical leader / TL |
| Prompt / issue | PO-02b (`docs/prompts/PO-02b-ui-spec-decisions.md`) |
| Branch | `docs/po-02-ui-spec` |
| PR | [#8](https://github.com/duong2604/githive/pull/8) |
| Last commit | `docs: record PO-02b` |
| State | waiting for Danny to merge (PR #8 goes ready when CI is green) |
| Last updated | 2026-11-14 |

Previous sessions: PO-01, PR #7 merged (563cf09); PO-02 on this same PR.

## Session goal

Apply the PO decisions UI-D1…UI-D10: UI spec Approved (E1…E12), BD-00 overview rev 4 (O1…O5), CLAUDE.md §4 without "sync strategy".

## Done

- Prompt saved (R3); UI spec Approved with §19 and revision history; BD-00 overview rev 4; CLAUDE.md §4 and README `init` row.
- PROGRESS, prompt index, report `docs/reports/PO-02b-2026-11-14.md`.

## In progress

Nothing.

## Next steps

1. If CI on PR #8 fails: fix on this branch, push, update this file.
2. After Danny merges PR #8, next TL session = P05 (PO-02-N4), on a new branch from `origin/main`:
   1. DD-00 rev 4: lock-info sharing on Windows (retry, warn), `--json` on flag errors in any order, timing tolerances in TC-20/21/28, `Prompter` with `Option`/default/validate and Esc/`q` → `interrupted` 130 (UI-D6), `Mode.Unicode` (false when not a TTY, UI-D5) and `StderrTTY`, `Printer.Success`/`Warn` format, help template, cobra `completion` kept (UI-D7), gate bypass list (UI-D9), mintty tip (UI-D10), version from `debug.ReadBuildInfo`, owner TL. Keep `HuhPrompter` out of DD-00a (moves to DD-01a).
   2. DD-03 rev 4: `Spec.Command`, export `Resolver.Current`, picker options with labels, `NULL` for a missing base (UI-D2), owner TL.
   3. Update P05 (proposals in `docs/reports/PO-01-2026-11-14.md` b), create the 5 issues, write P06 (review criteria include PO-02-N1: `version --json` uses the common envelope).
   4. Apply PO answers to Q9, Q10 (`docs/reports/PO-02b-2026-11-14.md` §6) when they arrive.

## Commands to resume

```
git fetch origin && git switch docs/po-02-ui-spec && git pull
go run ./tools/docs-check && go vet ./... && go test ./...
gh pr checks 8
```

## Open questions

- PO: Q9 (wizard "not set" vs `NULL`), Q10 (BD-00 overview §7 `init` row still lists sync strategy).

## Decisions made in this session

- Lost backticks and `*` in the paste read by intent (E1 backticks, E7 `'kirindo@*'`); raw paste kept in the prompt file.
- UI spec §19 placed after §18 so section numbers stay in order.
- README `init` row fixed together with CLAUDE.md §4.
