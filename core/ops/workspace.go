package ops

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name string
	Path string
}

func ReadFile(path string) ([]byte, error) {
	rel, err := filepath.Rel(".", path)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(rel)
}

func ReadDir(root string) ([]File, error) {
	var files []File

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		name := d.Name()

		if strings.HasPrefix(name, ".") {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		files = append(files, File{
			Path: rel,
			Name: name,
		})

		return nil
	})

	return files, err
}
