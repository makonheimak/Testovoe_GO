package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var ErrEmptyConfig = errors.New("empty config")

func Decode(data []byte, source string) (Document, error) {
	return DecodeWithFormat(data, source, DetectFormat(source, data))
}

func DecodeWithFormat(data []byte, source string, format Format) (Document, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return Document{}, ErrEmptyConfig
	}

	var root any
	switch format {
	case FormatJSON:
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		if err := decoder.Decode(&root); err != nil {
			return Document{}, fmt.Errorf("decode json: %w", err)
		}
	case FormatYAML:
		if err := yaml.Unmarshal(data, &root); err != nil {
			return Document{}, fmt.Errorf("decode yaml: %w", err)
		}
	default:
		return Document{}, fmt.Errorf("unsupported config format for %q", source)
	}

	return Document{
		Source: source,
		Format: format,
		Root:   normalize(root),
	}, nil
}

func normalize(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, child := range typed {
			out[key] = normalize(child)
		}
		return out
	case map[any]any:
		out := make(map[string]any, len(typed))
		for key, child := range typed {
			out[fmt.Sprint(key)] = normalize(child)
		}
		return out
	case []any:
		out := make([]any, len(typed))
		for i, child := range typed {
			out[i] = normalize(child)
		}
		return out
	default:
		return typed
	}
}
