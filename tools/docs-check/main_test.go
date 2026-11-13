package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validFM = "---\nid: BD-04\ntitle: Spawn\nstatus: Draft   # comment\nreq: [REQ-04]\ndepends_on: []\nowner: Danny\nlast_updated: 2026-11-13\n---\n\n# BD-04\n"

const validPrompt = "---\nid: P01b\ntitle: Apply review\nstep: \"1b\"\nstatus: Ready\nbranch: docs/x\npr: \"\"\ndepends_on: [P01]\ndecisions: []\nlast_updated: 2026-11-13\n---\n"

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
			"docs/templates/BD-template.md":       "---\nid: BD-xx\nlast_updated: YYYY-MM-DD\n---\n",
			"docs/assets/flow.svg":                "<svg/>",
			"docs/detail-design/.gitkeep":         "",
			"docs/adr/0001-record-decisions.md":   validFM,
			"docs/adr/0002-rejected.md":           strings.Replace(validFM, "status: Draft", "status: Rejected", 1),
			"docs/prompts/P01b-apply-review.md":   validPrompt,
			"docs/prompts/README.md":              "---\nid: PROMPTS-INDEX\nlast_updated: 2026-11-13\n---\n",
			"docs/PROGRESS.md":                    "# no front matter\n",
			"docs/reports/P04-2026-11-13.md":      "# report, no front matter\n",
			"docs/reports/P04b-2026-11-13-2.md":   "# second session\n",
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
		{"b: Rejected outside ADR", map[string]string{"docs/basic-design/BD-04-spawn.md": strings.Replace(validFM, "status: Draft", "status: Rejected", 1)}, []string{`status "Rejected"`}},
		// Rule c
		{"c: uppercase slug", map[string]string{"docs/basic-design/BD-04-Spawn.md": validFM}, []string{"file name must match"}},
		{"c: one-digit number", map[string]string{"docs/detail-design/DD-4-spawn.md": validFM}, []string{"file name must match"}},
		// Rule d
		{"d: CRLF", map[string]string{"docs/requirements/x.md": "a\r\nb\r\n"}, []string{"CRLF"}},
		{"d: UTF-16 root README", map[string]string{"README.md": "\xff\xfe#\x00 \x00x\x00", "docs/x.md": "ok\n"}, []string{"not valid UTF-8"}},
		// Rule e
		{"e: dd-mm-yyyy date", map[string]string{"docs/basic-design/BD-04-spawn.md": strings.Replace(validFM, "2026-11-13", "13-11-2026", 1)}, []string{"not a valid YYYY-MM-DD"}},
		{"e: impossible date", map[string]string{"docs/adr/0003-x.md": strings.Replace(validFM, "2026-11-13", "2026-02-30", 1)}, []string{"not a valid YYYY-MM-DD"}},
		{"e: date in any front matter", map[string]string{"docs/PROGRESS.md": "---\nid: PROGRESS\nlast_updated: 2026/11/13\n---\n"}, []string{"not a valid YYYY-MM-DD"}},
		// Rule f
		{"f: prompt name", map[string]string{"docs/prompts/KICKOFF_PROMPT.md": validPrompt}, []string{"prompt file name must match"}},
		{"f: prompt missing field", map[string]string{"docs/prompts/P02-bootstrap.md": strings.Replace(validPrompt, "decisions: []\n", "", 1)}, []string{"missing field decisions"}},
		{"f: prompt status", map[string]string{"docs/prompts/P02-bootstrap.md": strings.Replace(validPrompt, "status: Ready", "status: Approved", 1)}, []string{`status "Approved"`}},
		{"f: prompt without front matter", map[string]string{"docs/prompts/P03-issues.md": "# P03\n"}, []string{"missing YAML front matter"}},
		// Rule g
		{"g: report without date", map[string]string{"docs/reports/P04-report.md": "# x\n"}, []string{"report file name must match"}},
		{"g: report bad prompt ID", map[string]string{"docs/reports/P4-2026-11-13.md": "# x\n"}, []string{"report file name must match"}},
		{"g: report dd-mm-yyyy", map[string]string{"docs/reports/P04-13-11-2026.md": "# x\n"}, []string{"report file name must match"}},
		{"g: report impossible date", map[string]string{"docs/reports/P04b-2026-02-30.md": "# x\n"}, []string{"not a valid YYYY-MM-DD"}},
		{"g: report bad suffix", map[string]string{"docs/reports/P04-2026-11-13-b.md": "# x\n"}, []string{"report file name must match"}},
		{"g: README in reports", map[string]string{"docs/reports/README.md": "# x\n"}, []string{"report file name must match"}},
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
