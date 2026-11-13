# GitHive

GitHive gives each parallel AI coding agent (or each human task) its own git worktree and branch, always created from the latest commit of a project's base branch, and cleans up safely afterwards.

**Status:** pre-alpha. Only `githive version` exists; the commands below are planned and are being designed in [`docs/basic-design`](docs/basic-design/README.md).

## Why

- **Multi-project:** one workspace, many repos, addressed as `project@task`.
- **Base-branch discipline:** each project has its own base branch (`dev`, `develop`, …), fetched before every spawn; new branches start from `origin/<base>`, never from a stale local branch.
- **Safe batch cleanup:** `retire` refuses to delete work that is not committed or pushed, and recognises squash-merged branches.
- **Windows/WSL first-class:** path rules, long paths, file locks and WSL-specific warnings.
- **Forge-agnostic:** works with GitLab, GitHub or self-hosted remotes.

## Planned commands

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

A target is `project@task`, or just `task` when you run the command inside a registered repo or one of its worktrees. `'project@*'` (quoted in POSIX shells) selects every task of one project for `swarm`, `sync` and `retire`; `swarm --all` lists every registered project. There is no `-p/--project` flag: the project comes from the target or the current directory. Task names are lowercase: `^[a-z0-9][a-z0-9._-]{0,49}$`. See [BD-03](docs/basic-design/BD-03-target.md). Every command accepts `--json` and `--yes`. Release packages also install `git-hive`, so `git hive <cmd>` works.

## Build

Requires Go (version in `go.mod`) and git 2.38 or newer.

```
go build ./...
go vet ./...
golangci-lint run
go test ./...
go run ./cmd/githive --help
```

## Documentation

- Requirements: [docs/requirements](docs/requirements/README.md)
- Basic design: [docs/basic-design](docs/basic-design/README.md)
- Decisions: [docs/adr](docs/adr/)
- Progress: [docs/PROGRESS.md](docs/PROGRESS.md)
- Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)

## License

[MIT](LICENSE)
