// Command docs-check validates the docs rules (docs/templates/, CLAUDE.md §3). Run from the repo root:
//
//	go run ./tools/docs-check [repo-root]
//
// Rules:
//
//	a) docs/ holds only .md, plus .png/.svg in docs/assets/ and .gitkeep (P01b-D3);
//	b) BD/DD/ADR front matter: required fields, allowed status (Rejected: ADR only, P01b-D1);
//	c) BD/DD file names match nameRe, except BD-00-overview.md;
//	d) .md files in docs/ and at the repository root are UTF-8 without CRLF (P01b-D3);
//	e) every front matter last_updated is a valid YYYY-MM-DD date (P01b-D5);
//	f) docs/prompts/ files except README.md match promptNameRe and have prompt front matter (P01b-D7);
//	g) docs/reports/ files match reportNameRe with a valid date; no front matter required (P04b-D10).
//
// Size (P02-D3): there is no fixed line limit. When this file grows past about
// 250 lines, split it into several files in this package (e.g. rules.go,
// frontmatter.go).
package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

type schema struct {
	required []string
	statuses []string
}

var (
	nameRe       = regexp.MustCompile(`^(BD|DD)-\d{2}-[a-z0-9-]+\.md$`)
	promptNameRe = regexp.MustCompile(`^P\d{2}[a-z]?-[a-z0-9-]+\.md$`)
	reportNameRe = regexp.MustCompile(`^P\d{2}[a-z]?-(\d{4}-\d{2}-\d{2})(-\d+)?\.md$`)
	dateRe       = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	nonEmpty     = map[string]bool{"id": true, "title": true, "status": true}
	docFields    = []string{"id", "title", "status", "req", "depends_on", "owner", "last_updated"}
	designDoc    = schema{docFields, []string{"Draft", "In review", "Approved", "Superseded"}}
	adrDoc       = schema{docFields, []string{"Draft", "In review", "Approved", "Superseded", "Rejected"}}
	promptDoc    = schema{[]string{"id", "title", "step", "status", "branch", "pr", "depends_on", "decisions", "last_updated"},
		[]string{"Ready", "Running", "Done", "Superseded"}}
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	problems, err := check(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "docs-check:", err)
		os.Exit(1)
	}
	for _, p := range problems {
		fmt.Println(p)
	}
	if len(problems) > 0 {
		fmt.Printf("docs-check: %d problem(s)\n", len(problems))
		os.Exit(1)
	}
	fmt.Println("docs-check: ok")
}

// check returns one message per rule violation. Paths in messages use "/".
func check(root string) ([]string, error) {
	var problems []string
	add := func(rel, format string, args ...any) {
		problems = append(problems, rel+": "+fmt.Sprintf(format, args...))
	}

	rootMD, err := filepath.Glob(filepath.Join(root, "*.md"))
	if err != nil {
		return nil, err
	}
	for _, p := range rootMD {
		if _, err := readMD(p, filepath.Base(p), add); err != nil {
			return nil, err
		}
	}

	err = filepath.WalkDir(filepath.Join(root, "docs"), func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		r, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		rel := filepath.ToSlash(r)
		name := d.Name()
		ext := strings.ToLower(path.Ext(name))

		switch { // rule a
		case ext == ".md":
		case name == ".gitkeep":
			return nil
		case strings.HasPrefix(rel, "docs/assets/") && (ext == ".png" || ext == ".svg"):
			return nil
		default:
			add(rel, "only .md files are allowed in docs/ (images: .png/.svg in docs/assets/)")
			return nil
		}
		text, err := readMD(p, rel, add)
		if err != nil || strings.HasPrefix(rel, "docs/templates/") {
			return err
		}

		isBDDD := strings.HasPrefix(name, "BD-") || strings.HasPrefix(name, "DD-")
		isADR := strings.HasPrefix(rel, "docs/adr/") && name != "README.md"
		isPrompt := path.Dir(rel) == "docs/prompts" && name != "README.md"
		if isBDDD && name != "BD-00-overview.md" && !nameRe.MatchString(name) {
			add(rel, "file name must match %s", nameRe)
		}
		if isPrompt && !promptNameRe.MatchString(name) {
			add(rel, "prompt file name must match %s", promptNameRe)
		}
		if path.Dir(rel) == "docs/reports" { // rule g
			m := reportNameRe.FindStringSubmatch(name)
			if m == nil {
				add(rel, "report file name must match %s", reportNameRe)
			} else if _, err := time.Parse("2006-01-02", m[1]); err != nil {
				add(rel, "report date %q is not a valid YYYY-MM-DD date", m[1])
			}
		}
		var msgs []string
		switch {
		case isPrompt:
			msgs = checkFrontMatter(text, &promptDoc)
		case isADR:
			msgs = checkFrontMatter(text, &adrDoc)
		case isBDDD:
			msgs = checkFrontMatter(text, &designDoc)
		default:
			msgs = checkFrontMatter(text, nil) // rule e only, when front matter exists
		}
		for _, m := range msgs {
			add(rel, "%s", m)
		}
		return nil
	})
	return problems, err
}

// readMD implements rule d and returns the text with LF line endings.
func readMD(p, rel string, add func(string, string, ...any)) (string, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(data) {
		add(rel, "not valid UTF-8")
	} else if bytes.Contains(data, []byte("\r\n")) {
		add(rel, "has CRLF line endings (use LF)")
	}
	return strings.ReplaceAll(string(data), "\r\n", "\n"), nil
}

// checkFrontMatter implements rules b, e and the front matter part of f.
// With a nil schema, a document without front matter is accepted.
func checkFrontMatter(text string, s *schema) []string {
	if !strings.HasPrefix(text, "---\n") {
		if s == nil {
			return nil
		}
		return []string{"missing YAML front matter"}
	}
	block, _, ok := strings.Cut(text[4:], "\n---\n")
	if !ok {
		return []string{"front matter is not closed with ---"}
	}
	fields := map[string]string{}
	for _, line := range strings.Split(block, "\n") {
		key, val, ok := strings.Cut(line, ":")
		if !ok || strings.HasPrefix(line, " ") || strings.HasPrefix(line, "#") {
			continue
		}
		if i := strings.Index(val, " #"); i >= 0 {
			val = val[:i]
		}
		fields[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(val), `"'`)
	}
	var msgs []string
	if d, ok := fields["last_updated"]; ok {
		if _, err := time.Parse("2006-01-02", d); !dateRe.MatchString(d) || err != nil {
			msgs = append(msgs, fmt.Sprintf("last_updated %q is not a valid YYYY-MM-DD date", d))
		}
	}
	if s == nil {
		return msgs
	}
	for _, f := range s.required {
		v, ok := fields[f]
		switch {
		case !ok:
			msgs = append(msgs, "front matter missing field "+f)
		case v == "" && nonEmpty[f]:
			msgs = append(msgs, "front matter field "+f+" is empty")
		}
	}
	if st := fields["status"]; st != "" && !slices.Contains(s.statuses, st) {
		msgs = append(msgs, fmt.Sprintf("status %q is not one of %s", st, strings.Join(s.statuses, ", ")))
	}
	return msgs
}
