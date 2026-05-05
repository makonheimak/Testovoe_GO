package config

import (
	"bytes"
	"path/filepath"
	"strings"
)

type Format string

const (
	FormatUnknown Format = "unknown"
	FormatJSON    Format = "json"
	FormatYAML    Format = "yaml"
)

func DetectFormat(source string, data []byte) Format {
	if format := FormatFromFilename(source); format != FormatUnknown {
		return format
	}

	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return FormatUnknown
	}

	switch trimmed[0] {
	case '{', '[':
		return FormatJSON
	default:
		return FormatYAML
	}
}

func FormatFromFilename(source string) Format {
	switch strings.ToLower(filepath.Ext(source)) {
	case ".json":
		return FormatJSON
	case ".yaml", ".yml":
		return FormatYAML
	default:
		return FormatUnknown
	}
}
