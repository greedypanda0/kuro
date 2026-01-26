package config

import (
	"os"
	"path/filepath"
)

const DatabasePath = ".kuro/kuro.db"
const IgnorePath = ".kuro/.kuroignore"
const RepoDir = ".kuro"

func RepoRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := cwd
	for {
		if _, err := os.Stat(filepath.Join(current, RepoDir)); err == nil {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", os.ErrNotExist
		}
		current = parent
	}
}

func DatabasePathFor(root string) string {
	return filepath.Join(root, DatabasePath)
}

func IgnorePathFor(root string) string {
	return filepath.Join(root, IgnorePath)
}
