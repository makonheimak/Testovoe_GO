package dirscan

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"

	"config-audit/internal/config"
)

func FindConfigs(root string) ([]string, error) {
	var paths []string

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if config.FormatFromFilename(path) == config.FormatUnknown {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan configs in %q: %w", root, err)
	}

	sort.Strings(paths)
	return paths, nil
}
