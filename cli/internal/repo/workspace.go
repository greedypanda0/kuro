package repo

import (
	"cli/internal/config"
	coredb "core/db"
	"core/ops"
	"os"
	"path/filepath"
	"strings"
)

func ResetWorkspace(root string, db coredb.DBTX, snapshotHash string) error {
	snapshotFiles, err := coredb.ListSnapshotFiles(db, snapshotHash)
	if err != nil {
		return err
	}

	expected := make(map[string]string, len(snapshotFiles))
	for _, file := range snapshotFiles {
		expected[file.Path] = file.ObjectHash
	}

	currentFiles, err := ops.ReadDir(root)
	if err != nil {
		return err
	}

	for _, file := range currentFiles {
		if isKuroPath(file.Path) {
			continue
		}

		if _, ok := expected[file.Path]; ok {
			continue
		}
		absPath := filepath.Clean(filepath.Join(root, filepath.FromSlash(file.Path)))

		info, err := os.Lstat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}

		if info.IsDir() {
			if err := os.RemoveAll(absPath); err != nil {
				return err
			}
		} else {
			if err := os.Remove(absPath); err != nil {
				return err
			}
		}
	}

	for path, objectHash := range expected {
		obj, err := coredb.GetObject(db, objectHash)
		if err != nil {
			return err
		}

		absPath := filepath.Clean(filepath.Join(root, filepath.FromSlash(path)))
		if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(absPath, obj.Content, 0o644); err != nil {
			return err
		}
	}

	return nil
}

func isKuroPath(relPath string) bool {
	parts := strings.Split(filepath.Clean(relPath), string(os.PathSeparator))
	return len(parts) > 0 && parts[0] == config.RepoDir
}
