// Command docs-check validates the documentation rules described in
// docs/templates/ and CLAUDE.md §3. Run it from the repository root:
//
//	go run ./tools/docs-check [repo-root]
package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	nameRe         = regexp.MustCompile(`^(BD|DD)-\d{2}-[a-z0-9-]+\.md$`)
	requiredFields = []string{"id", "title", "status", "req", "depends_on", "owner", "last_updated"}
	nonEmpty       = map[string]bool{"id": true, "title": true, "status": true}
	statuses       = map[string]bool{"Draft": true, "In review": true, "Approved": true, "Superseded": true}
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

	// Rule d also covers the Markdown files at the repository root.
	rootMD, err := filepath.Glob(filepath.Join(root, "*.md"))
	if err != nil {
		return nil, err
	}
	for _, p := range rootMD {
		if err := checkEncoding(p, filepath.Base(p), add); err != nil {
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

		// Rule a: only Markdown in docs/, images only in docs/assets/.
		switch {
		case ext == ".md":
		case name == ".gitkeep":
		case strings.HasPrefix(rel, "docs/assets/") && (ext == ".png" || ext == ".svg"):
			return nil
		default:
			add(rel, "only .md files are allowed in docs/ (images: .png/.svg in docs/assets/)")
			return nil
		}
		if ext != ".md" {
			return nil
		}
		if err := checkEncoding(p, rel, add); err != nil {
			return err
		}
		if strings.HasPrefix(rel, "docs/templates/") {
			return nil
		}

		isBDDD := strings.HasPrefix(name, "BD-") || strings.HasPrefix(name, "DD-")
		isADR := strings.HasPrefix(rel, "docs/adr/") && name != "README.md"
		// Rule c: file name pattern.
		if isBDDD && name != "BD-00-overview.md" && !nameRe.MatchString(name) {
			add(rel, "file name must match %s", nameRe)
		}
		// Rule b: front matter.
		if isBDDD || isADR {
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			for _, msg := range checkFrontMatter(string(data)) {
				add(rel, "%s", msg)
			}
		}
		return nil
	})
	return problems, err
}

// checkEncoding implements rule d: UTF-8 without CRLF line endings.
func checkEncoding(p, rel string, add func(string, string, ...any)) error {
	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	if !utf8.Valid(data) {
		add(rel, "not valid UTF-8")
	} else if bytes.Contains(data, []byte("\r\n")) {
		add(rel, "has CRLF line endings (use LF)")
	}
	return nil
}

// checkFrontMatter implements rule b on a document's text.
func checkFrontMatter(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
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
		fields[strings.TrimSpace(key)] = strings.TrimSpace(val)
	}
	var msgs []string
	for _, f := range requiredFields {
		v, ok := fields[f]
		switch {
		case !ok:
			msgs = append(msgs, "front matter missing field "+f)
		case v == "" && nonEmpty[f]:
			msgs = append(msgs, "front matter field "+f+" is empty")
		}
	}
	if s, ok := fields["status"]; ok && s != "" && !statuses[s] {
		msgs = append(msgs, fmt.Sprintf("status %q is not one of Draft, In review, Approved, Superseded", s))
	}
	return msgs
}
