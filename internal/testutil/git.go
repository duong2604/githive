// Package testutil builds git fixtures for integration tests: a bare "remote"
// with main and dev branches plus a clone, isolated from the developer's git
// configuration (CLAUDE.md §10).
package testutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Fixture is a bare remote and a clone of it, both inside t.TempDir().
type Fixture struct {
	Remote string // bare repository, used as "origin"
	Clone  string // working clone with main checked out
}

// Isolate points git at an empty home and global config and fixes the author
// and committer, so tests never read or write the developer's real config.
// It uses t.Setenv, so the calling test must not be parallel.
func Isolate(t testing.TB) {
	t.Helper()
	home := t.TempDir()
	for k, v := range map[string]string{
		"HOME":                home,
		"USERPROFILE":         home,
		"XDG_CONFIG_HOME":     filepath.Join(home, ".config"),
		"GIT_CONFIG_GLOBAL":   filepath.Join(home, ".gitconfig"),
		"GIT_CONFIG_NOSYSTEM": "1",
		"GIT_AUTHOR_NAME":     "GitHive Test",
		"GIT_AUTHOR_EMAIL":    "test@githive.invalid",
		"GIT_AUTHOR_DATE":     "2026-01-01T00:00:00Z",
		"GIT_COMMITTER_NAME":  "GitHive Test",
		"GIT_COMMITTER_EMAIL": "test@githive.invalid",
		"GIT_COMMITTER_DATE":  "2026-01-01T00:00:00Z",
		"GIT_TERMINAL_PROMPT": "0",
		"LC_ALL":              "C",
	} {
		t.Setenv(k, v)
	}
}

// NewFixture isolates git (see Isolate) and creates a bare remote whose main
// branch has one commit and whose dev branch has one more, then clones it.
func NewFixture(t testing.TB) Fixture {
	t.Helper()
	Isolate(t)
	root := t.TempDir()
	f := Fixture{Remote: filepath.Join(root, "remote.git"), Clone: filepath.Join(root, "clone")}
	seed := filepath.Join(root, "seed")

	Git(t, root, "init", "--bare", "--initial-branch=main", f.Remote)
	Git(t, root, "init", "--initial-branch=main", seed)
	WriteFile(t, filepath.Join(seed, "README.md"), "seed\n")
	Git(t, seed, "add", "README.md")
	Git(t, seed, "commit", "-m", "initial commit")
	Git(t, seed, "switch", "-c", "dev")
	WriteFile(t, filepath.Join(seed, "dev.txt"), "dev\n")
	Git(t, seed, "add", "dev.txt")
	Git(t, seed, "commit", "-m", "dev commit")
	Git(t, seed, "push", f.Remote, "main", "dev")
	Git(t, root, "clone", f.Remote, f.Clone)
	return f
}

// Git runs git with args in dir and returns trimmed stdout; it fails the test
// on error. Arguments are passed as a slice, never through a shell.
func Git(t testing.TB, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		var stderr string
		if ee, ok := err.(*exec.ExitError); ok {
			stderr = string(ee.Stderr)
		}
		t.Fatalf("git %s (in %s): %v\n%s", strings.Join(args, " "), dir, err, stderr)
	}
	return strings.TrimSpace(string(out))
}

// WriteFile writes content to path, failing the test on error.
func WriteFile(t testing.TB, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
