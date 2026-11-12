package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	info := BuildInfo{Version: "1.2.3", Commit: "abc123", Date: "2026-11-13"}
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string
		wantStderr string
	}{
		{"version text", []string{"version"}, ExitOK, "githive 1.2.3 (commit abc123, built 2026-11-13)\n", ""},
		{"help", []string{"--help"}, ExitOK, "Usage:", ""},
		{"no args prints help", nil, ExitOK, "Usage:", ""},
		{"unknown command", []string{"spawnn"}, ExitUsage, "", `unknown command "spawnn"`},
		{"unknown flag", []string{"version", "--nope"}, ExitUsage, "", "unknown flag: --nope"},
		{"extra arg", []string{"version", "x"}, ExitUsage, "", "unknown command"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out, errOut bytes.Buffer
			code := Execute(info, tt.args, &out, &errOut)
			if code != tt.wantCode {
				t.Fatalf("exit code = %d, want %d (stderr %q)", code, tt.wantCode, errOut.String())
			}
			if !strings.Contains(out.String(), tt.wantStdout) {
				t.Errorf("stdout = %q, want it to contain %q", out.String(), tt.wantStdout)
			}
			if !strings.Contains(errOut.String(), tt.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", errOut.String(), tt.wantStderr)
			}
		})
	}
}

func TestVersionJSON(t *testing.T) {
	info := BuildInfo{Version: "1.2.3", Commit: "abc123", Date: "2026-11-13"}
	var out, errOut bytes.Buffer
	if code := Execute(info, []string{"version", "--json"}, &out, &errOut); code != ExitOK {
		t.Fatalf("exit code = %d, stderr %q", code, errOut.String())
	}
	var got BuildInfo
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON %q: %v", out.String(), err)
	}
	if got != info {
		t.Errorf("got %+v, want %+v", got, info)
	}
}
