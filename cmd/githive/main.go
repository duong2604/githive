// Command githive gives each parallel AI coding agent (or human task) its own
// git worktree and branch. See README.md.
package main

import (
	"os"

	"github.com/duong2604/githive/internal/cli"
)

// Set by GoReleaser through -ldflags "-X main.version=... -X main.commit=... -X main.date=...".
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	os.Exit(cli.Execute(cli.BuildInfo{Version: version, Commit: commit, Date: date}, os.Args[1:], os.Stdout, os.Stderr))
}
