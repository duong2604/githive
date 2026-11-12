package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFixture(t *testing.T) {
	f := NewFixture(t)

	if got := Git(t, f.Clone, "branch", "--show-current"); got != "main" {
		t.Errorf("clone is on %q, want main", got)
	}
	if got := Git(t, f.Clone, "rev-list", "--count", "origin/main..origin/dev"); got != "1" {
		t.Errorf("origin/dev is %s commits ahead of origin/main, want 1", got)
	}
	if cfg := Git(t, f.Clone, "config", "--list"); strings.Contains(cfg, "user.name=") {
		t.Errorf("user.name leaked from the developer's config:\n%s", cfg)
	}
}

func TestIsolateIgnoresRealConfig(t *testing.T) {
	Isolate(t)
	global := os.Getenv("GIT_CONFIG_GLOBAL")
	if filepath.Dir(global) != os.Getenv("HOME") {
		t.Fatalf("GIT_CONFIG_GLOBAL = %q, want it inside HOME %q", global, os.Getenv("HOME"))
	}
	dir := t.TempDir()
	Git(t, dir, "config", "--global", "githive.probe", "1")
	if got := Git(t, dir, "config", "--global", "--get", "githive.probe"); got != "1" {
		t.Errorf("isolated global config = %q, want 1", got)
	}
	if _, err := os.Stat(global); err != nil {
		t.Errorf("global config not written to the isolated file: %v", err)
	}
}
