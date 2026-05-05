package cli

import (
	"bytes"
	"errors"
	"testing"
)

func TestParseOptionsFilePath(t *testing.T) {
	options, err := ParseOptions([]string{"config.yaml"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ParseOptions returned error: %v", err)
	}

	if options.InputPath != "config.yaml" {
		t.Fatalf("expected input path config.yaml, got %q", options.InputPath)
	}
}

func TestParseOptionsShortSilent(t *testing.T) {
	options, err := ParseOptions([]string{"-s", "config.yaml"}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("ParseOptions returned error: %v", err)
	}

	if !options.Silent {
		t.Fatal("expected silent mode")
	}
}

func TestParseOptionsStdinRejectsFilePath(t *testing.T) {
	_, err := ParseOptions([]string{"--stdin", "config.yaml"}, &bytes.Buffer{})
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("expected usage error, got %v", err)
	}
}

func TestParseOptionsRejectsHTTPAndGRPCTogether(t *testing.T) {
	_, err := ParseOptions([]string{"--http", "--grpc"}, &bytes.Buffer{})
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("expected usage error, got %v", err)
	}
}
