package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validFM = "---\nid: BD-04\ntitle: Spawn\nstatus: Draft   # comment\nreq: [REQ-04]\ndepends_on: []\nowner: Danny\nlast_updated: 13-11-2026\n---\n\n# BD-04\n"

func TestCheck(t *testing.T) {
	tests := []struct {
		name  string
		files map[string]string
		want  []string // substrings, one per expected problem
	}{
		{"clean repo", map[string]string{
			"README.md":                           "# x\n",
			"docs/basic-design/BD-04-spawn.md":    validFM,
			"docs/basic-design/BD-00-overview.md": validFM,
			"docs/basic-design/README.md":         "# index\n",
			"docs/templates/BD-template.md":       "---\nid: BD-xx\n---\n",
			"docs/assets/flow.svg":                "<svg/>",
			"docs/detail-design/.gitkeep":         "",
			"docs/adr/0001-record-decisions.md":   validFM,
		}, nil},
		// Rule a
		{"a: docx in docs", map[string]string{"docs/requirements/x.docx": "x"}, []string{"only .md"}},
		{"a: png outside assets", map[string]string{"docs/basic-design/x.png": "x"}, []string{"only .md"}},
		{"a: pdf in assets", map[string]string{"docs/assets/x.pdf": "x"}, []string{"only .md"}},
		// Rule b
		{"b: no front matter", map[string]string{"docs/detail-design/DD-04-spawn.md": "# DD-04\n"}, []string{"missing YAML front matter"}},
		{"b: unclosed front matter", map[string]string{"docs/detail-design/DD-04-spawn.md": "---\nid: DD-04\n"}, []string{"not closed"}},
		{"b: missing field", map[string]string{"docs/detail-design/DD-04-spawn.md": strings.Replace(validFM, "owner: Danny\n", "", 1)}, []string{"missing field owner"}},
		{"b: empty title", map[string]string{"docs/detail-design/DD-04-spawn.md": strings.Replace(validFM, "title: Spawn", "title:", 1)}, []string{"title is empty"}},
		{"b: bad status", map[string]string{"docs/adr/0002-x.md": strings.Replace(validFM, "status: Draft", "status: Done", 1)}, []string{`status "Done"`}},
		// Rule c
		{"c: uppercase slug", map[string]string{"docs/basic-design/BD-04-Spawn.md": validFM}, []string{"file name must match"}},
		{"c: one-digit number", map[string]string{"docs/detail-design/DD-4-spawn.md": validFM}, []string{"file name must match"}},
		// Rule d
		{"d: CRLF", map[string]string{"docs/prompts/x.md": "a\r\nb\r\n"}, []string{"CRLF"}},
		{"d: UTF-16 root README", map[string]string{"README.md": "\xff\xfe#\x00 \x00x\x00", "docs/x.md": "ok\n"}, []string{"not valid UTF-8"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			for name, content := range tt.files {
				p := filepath.Join(root, filepath.FromSlash(name))
				if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got, err := check(root)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %d problems %q, want %d %q", len(got), got, len(tt.want), tt.want)
			}
			for i, w := range tt.want {
				if !strings.Contains(got[i], w) {
					t.Errorf("problem %d = %q, want it to contain %q", i, got[i], w)
				}
			}
		})
	}
}
