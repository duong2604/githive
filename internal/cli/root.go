// Package cli holds the cobra commands: it parses flags, calls services and
// renders output. It contains no business logic.
package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// Exit codes are a public API (CLAUDE.md §4).
const (
	ExitOK    = 0
	ExitError = 1
	ExitUsage = 2
)

// BuildInfo is injected at build time through -ldflags.
type BuildInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

type options struct {
	json bool
	yes  bool
}

// usageError marks bad arguments or flags (exit code 2).
type usageError struct{ err error }

func (e usageError) Error() string { return e.err.Error() }
func (e usageError) Unwrap() error { return e.err }

// Execute runs the CLI with args and returns the process exit code.
func Execute(info BuildInfo, args []string, stdout, stderr io.Writer) int {
	root := newRootCmd(info)
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	err := root.Execute()
	if err == nil {
		return ExitOK
	}
	_, _ = fmt.Fprintln(stderr, "githive:", err)
	var ue usageError
	if errors.As(err, &ue) {
		return ExitUsage
	}
	return ExitError
}

func newRootCmd(info BuildInfo) *cobra.Command {
	opts := &options{}
	root := &cobra.Command{
		Use:           "githive",
		Short:         "One git worktree and branch per AI agent or task, from the latest base branch",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) > 0 {
				return usageError{fmt.Errorf("unknown command %q (see githive --help)", args[0])}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
	}
	root.PersistentFlags().BoolVar(&opts.json, "json", false, "print machine-readable JSON")
	root.PersistentFlags().BoolVar(&opts.yes, "yes", false, "assume yes; never prompt")
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error { return usageError{err} })
	root.AddCommand(newVersionCmd(info, opts))
	return root
}

// noArgs is cobra.NoArgs reported as a usage error (exit code 2).
func noArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.NoArgs(cmd, args); err != nil {
		return usageError{err}
	}
	return nil
}

func newVersionCmd(info BuildInfo, opts *options) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version, commit and build date",
		Args:  noArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if opts.json {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(info)
			}
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "githive %s (commit %s, built %s)\n", info.Version, info.Commit, info.Date)
			return err
		},
	}
}
