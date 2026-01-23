package internal

import (
	"os"
	"path/filepath"
	"strings"
)

func ReadKuroIgnore(repoRoot string) ([]string, error) {
	path := filepath.Join(repoRoot, ".kuro", ".kuroignore")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
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
