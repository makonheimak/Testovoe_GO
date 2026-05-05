package config

import "testing"

func TestDecodeJSON(t *testing.T) {
	doc, err := Decode([]byte(`{"log":{"level":"debug"}}`), "config.json")
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if doc.Format != FormatJSON {
		t.Fatalf("expected format %q, got %q", FormatJSON, doc.Format)
	}

	root, ok := doc.Root.(map[string]any)
	if !ok {
		t.Fatalf("expected root map, got %T", doc.Root)
	}

	log, ok := root["log"].(map[string]any)
	if !ok {
		t.Fatalf("expected log map, got %T", root["log"])
	}

	if log["level"] != "debug" {
		t.Fatalf("expected debug level, got %v", log["level"])
	}
}

func TestDecodeYAML(t *testing.T) {
	doc, err := Decode([]byte("storage:\n  digest-algorithm: MD5\n"), "config.yaml")
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if doc.Format != FormatYAML {
		t.Fatalf("expected format %q, got %q", FormatYAML, doc.Format)
	}

	var found bool
	doc.Walk(func(node Node) {
		if node.Path == "storage.digest-algorithm" && node.Value == "MD5" {
			found = true
		}
	})

	if !found {
		t.Fatal("expected to find storage.digest-algorithm node")
	}
}

func TestDecodeEmptyConfig(t *testing.T) {
	if _, err := Decode([]byte(" \n\t"), "config.yaml"); err == nil {
		t.Fatal("expected error for empty config")
	}
}
