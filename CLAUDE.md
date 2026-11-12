# CLAUDE.md — GitHive

Instructions for Claude Code working in this repository. Read this file fully at the start of every session.

## 1. Team and roles

| Role | Who | Responsibility |
| --- | --- | --- |
| Client / product owner | Danny (`duong2604`) | Owns scope and priorities, approves PRs, merges to `main`, decides anything in "Ask first" (§9) |
| Technical advisor | Claude (claude.ai project "Git worktree tool") | Owns the architecture in `docs/design/basic-design.md`; reviews ADRs |
| Senior developer | **You** (Claude Code) | Turns issues into tested, reviewable PRs; proposes design changes as ADRs; flags risks early |

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

- Architecture and flows: `docs/design/basic-design.md`. It was written while the tool was still called `wtx`. Read `wtx` as `githive`, `wtx.*` as `githive.*`, `WTX_*` as `GITHIVE_*`, and the old command names as in the §4 table.
- Decisions made after the design doc: `docs/adr/`.
- If they conflict, the newest accepted ADR wins.

## 4. CLI surface (decided)

Binary: `githive`. Release packaging also installs `git-hive`, so `git hive <cmd>` works through git's external-command mechanism. A short `ghive` alias is an **open question**: don't add it until Danny decides.

| Command | Aliases | Purpose |
| --- | --- | --- |
| `githive init` | | Setup wizard: projects root, worktrees root, branch prefix, sync strategy |
| `githive project add <path>` / `list` / `config <p> --base dev` | `hive` | Register repos (repos under projects root are auto-discovered); list shows the base branch, or `NULL` when unset |
| `githive spawn <target>` | `new`, `create` | Fetch base, create worktree + branch `<prefix><task>`, copy files, run hooks, print the path |
| `githive swarm [--all]` | `ls`, `list` | List worktrees: branch, ahead/behind base, dirty, last commit |
| `githive inspect <target>` | | Detailed status of one worktree |
| `githive sync <target> [--rebase\|--merge] [--dry-run]` | | Fetch base, predict conflicts (`git merge-tree`), then rebase or merge |
| `githive retire <targets...> [--merged] [--stale 14d] [-i]` | `rm`, `remove` | Safe removal of worktree + local branch |
| `githive path <target>` | | Print the path (for `cd`) |
| `githive doctor` | | Environment checks |
| `githive config get/set` | | Global config |

**Target grammar:** `project@task`. Inside a registered repo or one of its worktrees, `task` alone is enough. `project@*` (quote it in shells) or `-p project --all` selects every task. Task names must match `^[a-z0-9][a-z0-9._-]{0,49}$` (no `@`, no `/`) and also pass `git check-ref-format --branch`. The grammar must never need a leading `@`, because PowerShell treats a leading `@` as splatting.

Every command accepts `--json` and `--yes` and never blocks waiting for input when stdin is not a TTY. Exit codes: 0 ok, 1 error, 2 bad args/missing config, 3 sync conflict (aborted cleanly), 4 blocked by a safety check, 5 lock timeout. `--json` schemas and exit codes are a public API: changing them after v0.1.0 needs Danny's approval.

## 5. Key technical decisions

- **Go** (latest stable; pin it in `go.mod`), a single static binary built with GoReleaser.
- **Shell out to the system `git`.** No go-git. Always pass an argument slice (never `sh -c` with concatenated strings), set `LC_ALL=C`, and parse `--porcelain -z` output. Minimum git version is 2.38.
- **Git is the source of truth** for worktrees (`git worktree list --porcelain -z`). No state database.
- **Config, two layers:** global TOML at `%APPDATA%\githive\config.toml` (Windows) or `$XDG_CONFIG_HOME/githive/config.toml` (default `~/.config/...`); per-repo settings in `git config --local githive.*`. Write files atomically (temp file, then rename).
- **Precedence:** flag > `GITHIVE_*` env > repo git config > global config > defaults.
- **Locking:** the `githive.lock` file in `git rev-parse --git-common-dir`, held only around fetch / worktree add / remove / branch delete. 30 s timeout.
- **Paths:** worktrees live at `<worktrees_root>/<project>/<task>`; the default `worktrees_root` is `<projects_root>/worktrees`. Never put worktrees inside the repo. Use `filepath` everywhere and never hand-build separators.
- **Spawn:** `git fetch <remote> <base> --prune`, then `git worktree add --no-track -b <prefix><task> <path> <remote>/<base>`. Roll back on failure.
- **Retire:** use `git worktree remove` + `git branch -d`, never raw recursive deletes. Before removing, refuse when the worktree is the main one, dirty, has unpushed commits, or is locked. Detect "merged" by ancestry first; otherwise treat a remote branch that has disappeared as merged.

Default dependencies (anything else → ask first): `spf13/cobra`, `charmbracelet/huh`, `charmbracelet/lipgloss`, `pelletier/go-toml/v2`, `gofrs/flock`, `rogpeppe/go-internal/testscript` (tests), `stretchr/testify` (tests).

Layout:
```
cmd/githive/            main.go
internal/cli/           cobra commands only: parse flags, call services, render
internal/ui/            prompts, tables, JSON output
internal/service/       project, worktree, sync, cleanup (no terminal I/O)
internal/git/           runner, porcelain parsers, status
internal/config/        global TOML, repo git config, resolution
internal/lock/          file lock
internal/platform/      paths, WSL detection, Windows drives
internal/target/        project@task parser
internal/testutil/      temp remote/clone fixtures
docs/design/  docs/adr/  testdata/
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
- Open PRs as **draft** and mark them ready only after CI is green. The body includes `Closes #N`, what/why, how it was tested (which OS locally), CLI output samples, and anything unverified.
- Keep PRs under ~400 changed lines excluding generated/testdata files; split otherwise.
- **Never:** push to `main`, merge PRs, force-push to anything except your own branch (and then only `--force-with-lease`), create tags or releases, or change repo settings or branch protection.
- **Parallel agents on this repo:** each agent works in its own worktree on its own issue. Until GitHive's MVP works, use `claude --worktree <name>` or `git worktree add ../githive-wt/<task> -b feat/<issue>-<slug> origin/main`. Once the MVP is released, dogfood GitHive itself (base branch `main`). Never `checkout`/`switch` another branch inside a worktree that isn't yours.

## 8. Changing the design (ADRs)

Write `docs/adr/NNNN-short-title.md` (Context / Decision / Consequences / Status: Proposed) when you want to:
- deviate from `docs/design/basic-design.md` or an accepted ADR;
- change the command surface, flags, exit codes, JSON schema, or config keys;
- add a dependency outside the default list.

Open the ADR in its own PR and tell Danny in Vietnamese what you propose and why. He'll take it to the advisor. Don't build on a Proposed ADR until it's Accepted.

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
