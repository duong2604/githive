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
| `githive project add <path>` / `list` / `config <p> --base dev` | `hive` | Register repos; `list` shows each base branch, or `NULL` when unset |
| `githive spawn <target>` | `new`, `create` | Fetch base, create worktree + branch `<prefix><task>`, copy files, run hooks, print the path |
| `githive swarm [--all]` | `ls`, `list` | List worktrees: branch, ahead/behind base, dirty, last commit |
| `githive inspect <target>` | | Detailed status of one worktree |
| `githive sync <target> [--rebase\|--merge] [--dry-run]` | | Fetch base, predict conflicts, then rebase or merge |
| `githive retire <targets...> [--merged] [--stale 14d] [-i]` | `rm`, `remove` | Safe removal of worktree + local branch |
| `githive path <target>` | | Print the path (for `cd`) |
| `githive doctor` | | Environment checks |
| `githive config get/set` | | Global config |
| `githive version` | | Print version, commit and build date (available now) |

A target is `project@task`; inside a registered repo or one of its worktrees, `task` alone is enough. Every command accepts `--json` and `--yes`. Release packages also install `git-hive`, so `git hive <cmd>` works.

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
