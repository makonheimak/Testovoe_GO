package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzePathFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"log":{"level":"debug"}}`), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	findings, err := NewDefaultAnalyzer().AnalyzePath(context.Background(), path, false, false)
	if err != nil {
		t.Fatalf("AnalyzePath returned error: %v", err)
	}

	if len(findings) != 1 || findings[0].RuleID != "debug-logging" {
		t.Fatalf("expected debug-logging finding, got %#v", findings)
	}
}

func TestAnalyzePathDirectoryRequiresRecursiveFlag(t *testing.T) {
	_, err := NewDefaultAnalyzer().AnalyzePath(context.Background(), t.TempDir(), false, false)
	if err == nil {
		t.Fatal("expected error for directory without recursive flag")
	}
}

func TestAnalyzePathRecursiveDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "debug.json"), []byte(`{"debug":true}`), 0o600); err != nil {
		t.Fatalf("WriteFile debug.json returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "weak.yaml"), []byte("storage:\n  digest-algorithm: MD5\n"), 0o600); err != nil {
		t.Fatalf("WriteFile weak.yaml returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "notes.txt"), []byte("ignore me"), 0o600); err != nil {
		t.Fatalf("WriteFile notes.txt returned error: %v", err)
	}

	findings, err := NewDefaultAnalyzer().AnalyzePath(context.Background(), root, true, false)
	if err != nil {
		t.Fatalf("AnalyzePath returned error: %v", err)
	}

	if len(findings) != 2 {
		t.Fatalf("expected two findings, got %#v", findings)
	}
}
