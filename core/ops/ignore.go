package ops

import (
	"github.com/greedypanda0/kuro/core/errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ReadKuroIgnore(ignorePath string) ([]string, error) {
	path := filepath.Join(ignorePath)

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

func IsIgnored(p string, ignores []string) bool {
	target := normalizePath(p)
	for _, raw := range ignores {
		pat := normalizePattern(raw)
		if pat == "" {
			continue
		}

		if strings.HasSuffix(pat, "/") {
			dir := strings.TrimSuffix(pat, "/")
			if dir != "" && (target == dir || strings.HasPrefix(target, dir+"/")) {
				return true
			}
			continue
		}

		if strings.Contains(pat, "/") {
			if matchPath(target, pat) {
				return true
			}
			continue
		}

		if matchAnySegment(target, pat) {
			return true
		}
	}

	return false
}

func normalizePath(p string) string {
	p = filepath.ToSlash(filepath.Clean(p))
	if p == "." {
		return ""
	}
	return strings.TrimPrefix(p, "./")
}

func normalizePattern(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	p = strings.TrimPrefix(p, "./")
	return path.Clean(path.Clean(p))
}

func matchPath(target, pattern string) bool {
	ok, err := path.Match(pattern, target)
	if err != nil {
		return false
	}
	return ok
}

func matchAnySegment(target, pattern string) bool {
	if target == "" {
		return false
	}
	segments := strings.Split(target, "/")
	for _, seg := range segments {
		ok, err := path.Match(pattern, seg)
		if err == nil && ok {
			return true
		}
	}
	return false
}
