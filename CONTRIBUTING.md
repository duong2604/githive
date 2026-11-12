# Contributing

GitHive is developed from designs: requirement (`REQ-xx`) → basic design (`BD-xx`) → detail design (`DD-xx`) → issue → pull request. Code is written only for a DD whose status is `Approved`.

## Workflow

1. Pick an issue. Each issue names the DD it implements.
2. Branch from an up-to-date `origin/main`: `feat/<issue>-<slug>`, `fix/…`, `docs/…`, `chore/…` or `test/…`.
3. Use [Conventional Commits](https://www.conventionalcommits.org/): `docs:` for documentation only, `chore:` for configuration and tooling, `feat:` only for user-visible features.
4. Open a draft pull request with `Closes #N`, the DD it implements and which OS you tested on. Mark it ready when CI is green.

Keep pull requests under about 400 changed lines, excluding generated files and test data.

## Checks

```
go build ./...
go vet ./...
golangci-lint run
go test ./...                 # add -race on Linux/macOS
go run ./tools/docs-check     # when docs/ changed
```

Integration tests use real git in temporary directories through `internal/testutil`, isolated from your own git configuration.

## Documentation

`docs/` contains only Markdown. Create new documents from [`docs/templates/`](docs/templates/) and use ISO dates (`YYYY-MM-DD`).
