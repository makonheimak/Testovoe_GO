package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ReadFile(path string) ([]byte, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("read config file %q: %w", path, err)
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return data, path, nil
	}

	return data, absolutePath, nil
}

func ReadAll(reader io.Reader, source string) ([]byte, string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", fmt.Errorf("read config from %s: %w", source, err)
	}
	return data, source, nil
}
