package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunStdinReturnsFindingsExitCode(t *testing.T) {
	stdin := strings.NewReader(`{"log":{"level":"debug"}}`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--stdin"}, stdin, &stdout, &stderr)
	if code != ExitFindings {
		t.Fatalf("expected exit code %d, got %d; stderr=%s", ExitFindings, code, stderr.String())
	}

	if !strings.Contains(stdout.String(), "debug-logging") {
		t.Fatalf("expected debug-logging output, got %q", stdout.String())
	}
}

func TestRunSilentReturnsOKWhenFindingsExist(t *testing.T) {
	stdin := strings.NewReader(`{"log":{"level":"debug"}}`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--stdin", "--silent"}, stdin, &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("expected exit code %d, got %d; stderr=%s", ExitOK, code, stderr.String())
	}
}

func TestRunJSONOutput(t *testing.T) {
	stdin := strings.NewReader(`{"storage":{"digest-algorithm":"MD5"}}`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--stdin", "--silent", "--json"}, stdin, &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("expected exit code %d, got %d; stderr=%s", ExitOK, code, stderr.String())
	}

	if !strings.Contains(stdout.String(), `"rule_id": "weak-algorithm"`) {
		t.Fatalf("expected JSON finding output, got %q", stdout.String())
	}
}

func TestRunFilePathReturnsOKForSafeConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "safe.json")
	if err := os.WriteFile(path, []byte(`{"log":{"level":"info"},"tls":{"enabled":true}}`), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{path}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("expected exit code %d, got %d; stderr=%s", ExitOK, code, stderr.String())
	}

	if !strings.Contains(stdout.String(), "No issues found.") {
		t.Fatalf("expected no issues message, got %q", stdout.String())
	}
}

func TestRunRecursiveDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "debug.json"), []byte(`{"debug":true}`), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "weak.yaml"), []byte("storage:\n  digest-algorithm: MD5\n"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--recursive", "--silent", root}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("expected exit code %d, got %d; stderr=%s", ExitOK, code, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "debug-logging") || !strings.Contains(out, "weak-algorithm") {
		t.Fatalf("expected recursive findings, got %q", out)
	}
}

func TestRunInvalidConfigReturnsErrorExitCode(t *testing.T) {
	stdin := strings.NewReader(`{"log":`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"--stdin"}, stdin, &stdout, &stderr)
	if code != ExitError {
		t.Fatalf("expected exit code %d, got %d", ExitError, code)
	}

	if !strings.Contains(stderr.String(), "decode json") {
		t.Fatalf("expected decode error, got %q", stderr.String())
	}
}
