package ops

import (
	"core/errors"
	"os"
	"path/filepath"
	"strings"
)

func ReadKuroIgnore(ignorePath string) ([]string, error) {
	path := ignorePath

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.ErrIgnoreFileNotFound
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var patterns []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		patterns = append(patterns, line)
	}

	return patterns, nil
}

func IsIgnored(path string, ignores []string) bool {
	for _, ig := range ignores {
		rel, err := filepath.Rel(ig, path)
		if err != nil {
			continue
		}

		if rel == "." || !strings.HasPrefix(rel, "..") {
			return true
		}
	}
	return false
}
