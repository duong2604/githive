# CLAUDE.md — GitHive

Instructions for Claude Code working in this repository. Read this file fully at the start of every session.

## 1. Team and roles

| Role | Who | Responsibility |
| --- | --- | --- |
| Project owner | Danny (`duong2604`) | Answers business and root-requirement questions, approves and merges PRs, decides anything in "Ask first" (§9), relays prompts between the roles |
| Product owner (PO) | Claude on claude.ai | Owns requirements, product and UI specification, and acceptance criteria; sends work to the TL as `PO-xx` prompts |
| Technical leader (TL) | Claude Code, TL session | Owns the technical design (DD, ADR) and technical decisions; splits work and writes Dev prompts (`P05`, `P06`…); reviews every Dev PR before Danny merges; reports progress to the PO |
| Developer (Dev) | Claude Code, Dev sessions (`DEV-1`, `DEV-2`…) | Writes code and tests from the TL's prompts; reports to the TL |

Before PO-01 a single "Advisor" (Claude on claude.ai) covered both product and technical decisions. From PO-01 on, product belongs to the PO and technical design to the TL. Advisor decisions recorded before PO-01 remain valid. The working agreement is in §14; handoff rules in §15.

Act like a senior engineer on a small team: small PRs, tests first where practical, honest status reports. When the design is silent, make a sensible call and record it in the PR. When you want to *contradict* the design, write an ADR (§8) and stop for approval.

Language: talk to Danny in **Vietnamese**. Write everything that lands in the repo — code, comments, commit messages, PRs, issues, docs, ADRs — in **English** (this is an open-source project).

## 2. Product

GitHive is a cross-platform CLI (Windows, WSL, Linux, macOS) that gives each parallel AI coding agent (or each human task) its own git worktree and branch, always created from the latest commit of a project's base branch, and cleans up safely afterwards.

**Positioning — why GitHive exists next to Worktrunk and `claude --worktree`:**
1. **Multi-project:** one workspace, many repos, addressed as `project@task`.
2. **Base-branch discipline:** a per-project base branch (`dev`, `develop`, …) that is fetched before every spawn; new branches start from `origin/<base>`, never from a stale local branch.
3. **Safe batch cleanup:** `retire` refuses to delete work that isn't committed or pushed, and recognises squash-merged branches.
4. **Windows/WSL first-class:** path rules, long paths, file locks, and WSL-specific warnings.
5. **Forge-agnostic core:** works with GitLab, GitHub, or self-hosted remotes. Forge APIs are optional extras (v2).

Do not copy code from other worktree tools. Reading their docs for ideas is fine.

## 3. Source of truth

```
requirements:    docs/requirements/README.md (REQ-xx, NFR-xx)
product:         docs/product/ui-spec.md — user-facing CLI interface (PO)
architecture:    docs/basic-design/BD-00-overview.md + docs/basic-design/README.md
feature designs: docs/basic-design/BD-xx-*.md → docs/detail-design/DD-xx-*.md
decisions:       docs/adr/
progress:        docs/PROGRESS.md
```

- `BD-00-overview.md` uses the official GitHive names since revision 2 (P01c); the original `wtx` draft is kept only in git history.
- If a design document and an ADR conflict, the newest Approved ADR wins.

**Format rules:**
- `docs/` contains only Markdown following `docs/templates/` (images, if any, as `.png`/`.svg` in `docs/assets/`; prefer Mermaid diagrams).
- Every new document is created from a template in `docs/templates/`.
- Run `go run ./tools/docs-check` before opening any PR that changes docs.
- `docs-check` also allows `.gitkeep` files (to keep empty folders) and checks the encoding of the `*.md` files at the repository root (P01b-D3).

**Workflow:** REQ → BD → DD → issue → PR.
- Implement only from a DD with `status: Approved`.
- Every issue and PR names its DD (e.g. "Implements DD-04").
- If the DD is missing or ambiguous, ask Danny. Do not guess.

## 4. CLI surface (decided)

Binary: `githive`. Release packaging also installs `git-hive`, so `git hive <cmd>` works through git's external-command mechanism. A short `ghive` alias is an **open question**: don't add it until Danny decides.

| Command | Aliases | Purpose |
| --- | --- | --- |
| `githive init` | | Setup wizard: projects root, worktrees root, branch prefix, sync strategy |
| `githive project add <path>` / `list` / `config <p> --base dev` | `hive` | Register repos (repos under projects root are auto-discovered); list shows the base branch, or `NULL` when unset |
| `githive spawn <target>` | `new`, `create` | Fetch base, create worktree + branch `<prefix><task>`, copy files, run hooks, print the path |
| `githive swarm ['project@*' \| --all]` | `ls`, `list` | List worktrees: branch, ahead/behind base, dirty, last commit; `--all` lists every registered project (ADR-0003) |
| `githive inspect <target>` | | Detailed status of one worktree |
| `githive sync <targets...> [--rebase\|--merge] [--dry-run]` | | Fetch base, predict conflicts (`git merge-tree`), then rebase or merge |
| `githive retire <targets...> [--merged] [--stale 14d] [-i] [--force]` | `rm`, `remove` | Safe removal of worktree + local branch; `--force` skips the CL-03 safety checks |
| `githive path <target>` | | Print the path (for `cd`) |
| `githive shell-init` | | Print a shell function for quick `cd` into a worktree (v1) |
| `githive doctor` | | Environment checks |
| `githive config get/set` | | Global config |
| `githive version` | | Print version, commit and build date; `--json` supported (since P02) |

**Target grammar:** `project@task`. Inside a registered repo or one of its worktrees, `task` alone is enough. `'project@*'` (quote it in shells) selects every task of one project for `swarm`, `sync` and `retire`; `swarm --all` lists every registered project; `-p project --all` is not used (ADR-0003). Task names must match `^[a-z0-9][a-z0-9._-]{0,49}$` (no `@`, no `/`) and also pass `git check-ref-format --branch`. The grammar must never need a leading `@`, because PowerShell treats a leading `@` as splatting. `-p/--project` does not exist; the project comes from the target or the current directory (BD-03, P03-D1).

Global flags: --json, --yes/-y, --verbose/-v, --no-color, --help (BD-00 §3.1). With --json, stdout holds one envelope {ok, command, data, error, warnings} (BD-00 §3.3, ADR-0004).

Running `githive` with no arguments prints help until DD-01 designs the interactive menu (P01d-D3).

Every command accepts `--json` and `--yes` and never blocks waiting for input when stdin is not a TTY. Exit codes: 0 ok, 1 error, 2 bad args/missing config, 3 sync conflict (aborted cleanly), 4 blocked by a safety check, 5 lock timeout, 130 interrupted. `--json` schemas and exit codes are a public API: changing them after v0.1.0 needs Danny's approval.

## 5. Key technical decisions

- **Go** (latest stable; pin it in `go.mod`), a single static binary built with GoReleaser.
- **Shell out to the system `git`.** No go-git. Always pass an argument slice (never `sh -c` with concatenated strings), set `LC_ALL=C`, and parse `--porcelain -z` output. Minimum git version is 2.38.
- **Git is the source of truth** for worktrees (`git worktree list --porcelain -z`). No state database.
- **Config, two layers:** global TOML at `%APPDATA%\githive\config.toml` (Windows) or `$XDG_CONFIG_HOME/githive/config.toml` (default `~/.config/...`); per-repo settings in `git config --local githive.*`. Write files atomically (temp file, then rename). The global key `extra_projects` lists repos outside `projects_root` registered with `githive project add` (P01c-D3). Key naming (P01d-D5): TOML keys use snake_case (`post_create`, `copy_files`); git config keys use camelCase under `githive.*` (`githive.postCreate`, `githive.copyFiles`).
- **Precedence:** flag > `GITHIVE_*` env > repo git config > global config > defaults.
- **Locking:** the `githive.lock` file in `git rev-parse --git-common-dir`, held only around fetch / worktree add / remove / branch delete / `git config --local` writes. 30 s timeout, overridable with `GITHIVE_LOCK_TIMEOUT` (seconds). The holder's command, pid, host and start time go in a separate `githive.lock.info` file next to it (ADR-0004).
- **Paths:** worktrees live at `<worktrees_root>/<project>/<task>`; the default `worktrees_root` is `<projects_root>/worktrees`. Never put worktrees inside the repo. Use `filepath` everywhere and never hand-build separators.
- **Spawn:** `git fetch <remote> <base> --prune`, then `git worktree add --no-track -b <prefix><task> <path> <remote>/<base>`. Roll back on failure.
- **Retire:** use `git worktree remove` + `git branch -d`, never raw recursive deletes. Before removing, refuse when the worktree is the main one, dirty, has unpushed commits, or is locked. Detect "merged" by ancestry first; otherwise treat a remote branch that has disappeared as merged.

Default dependencies (anything else → ask first): `spf13/cobra`, `charmbracelet/huh`, `charmbracelet/lipgloss`, `pelletier/go-toml/v2`, `gofrs/flock`, `golang.org/x/term`, `rogpeppe/go-internal/testscript` (tests), `stretchr/testify` (tests).

Layout:
```
cmd/githive/            main.go
internal/cli/           cobra commands only: parse flags, call services, render
internal/apperr/        errors: Error type, code catalogue, exit codes
internal/ui/            prompts, tables, JSON output
internal/service/       project, worktree, sync, cleanup (no terminal I/O)
internal/git/           runner, porcelain parsers, status
internal/config/        global TOML, repo git config, resolution
internal/lock/          file lock
internal/platform/      paths, WSL detection, Windows drives
internal/target/        project@task parser
internal/testutil/      temp remote/clone fixtures
tools/docs-check/       documentation checker
docs/requirements/  docs/basic-design/  docs/detail-design/  docs/templates/  docs/adr/  docs/prompts/
testdata/
```

## 6. Development environment

Danny works on two machines. Both are first-class, and they sync **only through GitHub**:

| Machine | OS / shell | Clone location |
| --- | --- | --- |
| Personal PC | Windows (Claude Code native, Git Bash / PowerShell) | `D:\workspace\githive` (or `%USERPROFILE%\workspace\githive`) |
| Work laptop | WSL (Ubuntu) | `~/workspace/githive`, **not** under `/mnt/c` or `/mnt/d` |

Rules:
- At the start of a session, detect the environment (`go env GOOS`, `uname -a` when available, `$WSL_DISTRO_NAME`) and say which one you're on.
- Never operate on a checkout with the "other side's" git (Windows git on a WSL clone or the reverse). Worktree `.git` files contain absolute paths.
- Line endings are LF everywhere and enforced by `.gitattributes`. Don't change `core.autocrlf` in the user's global config.
- You can only run Windows-specific behaviour on the Windows machine and Unix behaviour on WSL. Rely on the CI matrix for the other side, and **say explicitly in the PR which OS you verified locally**.
- No Makefile (Windows has no `make`). Use these commands:
  ```
  go build ./...
  go vet ./...
  golangci-lint run
  go test ./...                 # add -race on Linux/macOS
  go test ./... -run Integration
  go run ./cmd/githive --help
  ```

## 7. GitHub workflow

- Repo: `github.com/duong2604/githive`. Default branch: `main`. Use the `gh` CLI for issues and PRs.
- Work only from an issue. One issue per PR. Branch from up-to-date `origin/main` as `feat/<issue>-<slug>`, `fix/…`, `docs/…`, `chore/…`, `test/…`.
- Use Conventional Commits (`feat(spawn): fetch base before creating worktree`).
- Commits that only touch documentation use `docs:`; configuration and tooling use `chore:`; `feat:` is only for user-visible features.
- Open PRs as **draft** and mark them ready only after CI is green. The body includes `Closes #N`, what/why, how it was tested (which OS locally), CLI output samples, and anything unverified.
- Keep PRs under ~400 changed lines excluding generated/testdata files; split otherwise.
- **Never:** push to `main`, merge PRs, force-push to anything except your own branch (and then only `--force-with-lease`), create tags or releases, or change repo settings or branch protection.
- **Parallel agents on this repo:** each agent works in its own worktree on its own issue. Until GitHive's MVP works, use `claude --worktree <name>` or `git worktree add ../githive-wt/<task> -b feat/<issue>-<slug> origin/main`. Once the MVP is released, dogfood GitHive itself (base branch `main`). Never `checkout`/`switch` another branch inside a worktree that isn't yours.

## 8. Changing the design (ADRs)

Write `docs/adr/NNNN-short-title.md` (Context / Decision / Consequences / status: In review) when you want to:
- deviate from `docs/basic-design/BD-00-overview.md` (or an approved BD/DD) or an Approved ADR;
- change the command surface, flags, exit codes, JSON schema, or config keys;
- add a dependency outside the default list.

Open the ADR in its own PR and tell Danny in Vietnamese what you propose and why. He'll take it to the advisor. Don't build on an ADR until it is Approved.

## 9. Ask Danny first

- Anything in §8.
- The license choice, the `ghive` alias, release/versioning.
- Deleting or rewriting files your task didn't create.
- Any command that touches the user's real global git config, real repos outside temp dirs, or the network beyond `git fetch`/`gh` for this repo.

## 10. Testing standards

- **Unit tests:** parsers (porcelain, target grammar, config resolution) are table-driven, with golden files under `testdata/`.
- **Integration tests:** use real git in `t.TempDir()` via `internal/testutil` (a bare "remote" with `main` and `dev`, plus a clone). Isolate from the developer's config with `HOME`/`USERPROFILE`, `XDG_CONFIG_HOME`, `GIT_CONFIG_GLOBAL`, `GIT_CONFIG_NOSYSTEM=1`, and a fixed author/committer.
- **CLI end-to-end:** `testscript` scripts in `testdata/script/*.txtar`.
- **Concurrency:** a test where several goroutines spawn in the same repo at once; all must succeed without ref-lock errors.
- **Windows:** tests must not assume `/`, case-sensitive paths, or that open files can be deleted.

## 11. Definition of done

- [ ] Acceptance criteria of the issue met.
- [ ] `go vet`, `golangci-lint run`, and `go test ./...` pass locally; CI is green on ubuntu, windows, and macos.
- [ ] New behaviour has tests; bug fixes have a regression test.
- [ ] `--help` text, README command table, and `CHANGELOG.md` (Unreleased) updated when user-facing.
- [ ] PR description states the OS verified locally and anything left unverified.
- [ ] Final report to Danny in Vietnamese: what changed, what was tested and where, risks, open questions.

## 12. Security basics

- Never interpolate user input into a shell. The only intentional shell execution is a user-configured `postCreate` hook, run with the worktree as cwd.
- Validate task and project names before building paths; refuse any path that resolves outside `worktrees_root`.
- Hooks come only from the user's global config or the repo's local git config. A repo-committed `.githive.toml` (v2) must require explicit user trust before its hooks run.
- Never print tokens or remote URLs that contain credentials (strip `user:pass@`).

## 13. Reporting & escalation

End every session with a report to Danny in Vietnamese, at most ~40 lines, using exactly these headings:

```
## Báo cáo — <prompt ID, e.g. P02> — <YYYY-MM-DD>
1. Kết quả: Done | Partial | Blocked — one sentence
2. Nhánh / PR / commit cuối
3. Đã làm (≤ 5 bullets)
4. Lệch so với prompt hoặc thiết kế — every deviation, however small: what and why ("Không có" if none)
5. Câu hỏi kỹ thuật → TL — each with your recommended option (in the TL's own reports: technical decisions taken)
6. Câu hỏi requirement/spec → PO / Danny — each with the interpretation you used meanwhile
7. Kiểm chứng — OS, commands run, results; what is NOT verified
8. Việc Danny cần làm — exact commands
```

A TL report sent to the PO adds a section `9. PO summary`: at most 10 lines, in business language — which features are done, where they match or deviate from the requirements, and what the PO or Danny must decide.

Routing (§14): technical questions (architecture, libraries, git behaviour, file formats, conventions) are decided by the TL (the advisor before PO-01); Danny only reviews. Product and UX questions (what the user sees, defaults, safety policy, naming) go to the PO; business and root-requirement questions go to Danny. When unsure which kind, put it under 5. Do not block on a question if a reversible default exists: use it, list it, continue.
Every prompt lives in docs/prompts/ with an ID (see docs/prompts/README.md). Decisions are cited by full ID (e.g. P01b-D5); requirement clarifications by CL-xx. Update docs/PROGRESS.md and docs/prompts/README.md before writing the report.

Before printing the report, save it as docs/reports/<prompt ID>-<YYYY-MM-DD>.md (add -2, -3… for more sessions of the same prompt on one day) and commit it with the session's other changes, so any machine can read it after git pull (P04b-D10).

## 14. Working agreement

Ownership (PO-01):

| Area | Owner |
| --- | --- |
| `docs/requirements/`, `docs/product/` | PO |
| `docs/basic-design/` (behaviour the user sees) | PO writes; TL reviews feasibility |
| `docs/detail-design/`, `docs/adr/` | TL |
| Code and tests | Dev |
| PR review | TL |
| Merge | Danny |

Flow: `PO-xx` → TL → (`Pxx` prompt) → Dev → Dev report → TL review → TL report with "PO summary" (§13) → Danny relays it to the PO.

Rules for keeping history:

- **R1.** Nobody downloads files and copies them over the repo. Every change is made by an agent on a branch, committed, and opened as a PR.
- **R2.** Content written by the PO arrives only through prompts. For an existing document the PO sends section-by-section edits; full text only for a new document.
- **R3.** The first commit of every `PO-xx` session saves the prompt verbatim as `docs/prompts/PO-<NN>-<slug>.md`: the text Danny pasted, plus prompt front matter with `id: PO-<NN>`, `from: PO`, `to: TL`.
- **R4.** Prompts the TL writes for Devs are created in the repo (`docs/prompts/P<NN>-<slug>.md`, `from: TL`, `to: Dev`) and committed through a PR. Danny opens the file and pastes it; no file copying.
- **R5.** Every session saves its report as `docs/reports/<ID>-<YYYY-MM-DD>.md` (P04b-D10), including `PO-xx` sessions.

Question routing: technical questions are decided by the TL; product or UX questions go to the PO; business or root-requirement questions go to Danny.

## 15. Handoff and takeover

**R6. Reports by milestone.**
- End of every prompt or session: a report as in R5.
- End of every issue: the Dev saves the issue report in that issue's own PR (`docs/reports/P<NN>-<YYYY-MM-DD>.md`).
- When every slice of a DD is merged (the feature is done): the TL writes a feature report `docs/reports/FEAT-DD-<NN>-<YYYY-MM-DD>.md` covering what the feature does against the REQ/BD, each acceptance criterion checked, deviations, and remaining work.

**R7. A handoff file is always ready,** so another agent can take over at once if the current one goes offline.
- Every agent keeps a handoff file on the branch it works on: `docs/handoff/TL.md`, `docs/handoff/DEV-1.md`, `DEV-2.md`… Each Dev running in parallel has its own ID; Danny writes it on the first line of the prompt.
- The file follows `docs/templates/HANDOFF-template.md`: role and agent ID, prompt or issue, branch, PR, last commit; session goal; done; in progress (exact file and function); numbered next steps; commands to resume (`git switch …`, `go test …`); open questions and to whom; decisions made in the session; last update time.
- Update and push the handoff file after every group of commits, before every long command or CI wait, and before ending the session. Never leave more than one group of commits without updating it.
- Right after the first commit of a session, open the PR as a draft, so `gh pr list` always shows who is working on what, on which branch.
- When the PR is merged, the handoff file goes into `main` with it and stays as that agent's last trace.

**R8. Taking over from another agent.**
- The new agent runs `gh pr list --state open`, `git switch <branch>`, `git pull`, then reads `docs/handoff/<ID>.md`, `CLAUDE.md` and the latest report in `docs/reports/`.
- It adds a line `Taken over by <new ID> at <time>` to the handoff file and continues from "Next steps".
- It does not redo work that is already committed.
