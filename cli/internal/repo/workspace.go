package repo

import (
	"cli/internal/config"
	coredb "core/db"
	"core/ops"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ResetWorkspace(root string, db coredb.DBTX, snapshotHash string) error {
	snapshotFiles, err := coredb.ListSnapshotFiles(db, snapshotHash)
	if err != nil {
		return err
	}

	expected := make(map[string]struct{}, len(snapshotFiles))
	for _, f := range snapshotFiles {
		expected[f.Path] = struct{}{}
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

		abs := filepath.Join(root, filepath.FromSlash(file.Path))
		if err := os.Remove(abs); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	for _, f := range snapshotFiles {
		obj, err := coredb.GetObject(db, f.ObjectHash)
		if err != nil {
			return err
		}

		abs := filepath.Join(root, filepath.FromSlash(f.Path))
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(abs, obj.Content, 0o644); err != nil {
			return err
		}
	}

	var dirs []string

	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() || path == root {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err == nil && isKuroPath(rel) {
			return filepath.SkipDir
		}

		dirs = append(dirs, path)
		return nil
	})
	if err != nil {
		return err
	}

	sort.Slice(dirs, func(i, j int) bool {
		return len(dirs[i]) > len(dirs[j])
	})

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		if len(entries) == 0 {
			_ = os.Remove(dir)
		}
	}

	return nil
}

func isKuroPath(relPath string) bool {
	parts := strings.Split(filepath.Clean(relPath), string(os.PathSeparator))
	return len(parts) > 0 && parts[0] == config.RepoDir
}
