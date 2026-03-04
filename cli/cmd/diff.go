package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/greedypanda0/kuro/cli/internal/config"
	"github.com/greedypanda0/kuro/cli/internal/ui"

	coredb "github.com/greedypanda0/kuro/core/db"
	coreerrors "github.com/greedypanda0/kuro/core/errors"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

var diffCommand = &cobra.Command{
	Use:          "diff",
	Short:        "Show staged file changes",
	Long:         "Show differences between staged files and the last commit",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := config.RepoRoot()
		if err != nil {
			ui.Println(ui.Error("Repository not initialized"))
			return err
		}

		db, err := coredb.OpenDB(config.DatabasePathFor(root))
		if err != nil {
			ui.Println(ui.Error("Failed to open repository"))
			return err
		}
		defer db.Close()

		head, err := coredb.GetConfig(db, "head")
		if err != nil {
			ui.Println(ui.Error("Failed to get HEAD"))
			return err
		}

		ref, err := coredb.GetRef(db, head)
		if err != nil {
			ui.Println(ui.Error("Failed to get ref"))
			return err
		}

		stageFiles, err := coredb.GetStageFiles(db)
		if err != nil {
			ui.Println(ui.Error("Failed to get staged files"))
			return err
		}
		if len(stageFiles) == 0 {
			ui.Println(ui.Simple("No files staged"))
			return nil
		}

		paths := make([]string, 0, len(stageFiles))
		stagedSet := make(map[string]struct{}, len(stageFiles))
		for _, f := range stageFiles {
			paths = append(paths, f.Path)
			stagedSet[f.Path] = struct{}{}
		}
		sort.Strings(paths)

		fileFlag, _ := cmd.Flags().GetString("file")
		if fileFlag != "" {
			rel, err := resolveDiffPath(root, fileFlag)
			if err != nil {
				ui.Println(ui.Error("Invalid file path"))
				return err
			}
			if _, ok := stagedSet[rel]; !ok {
				ui.Println(ui.Error("File is not staged"))
				return nil
			}
			paths = []string{rel}
		}

		var snapshotHash string
		if ref != nil && ref.SnapshotHash != nil {
			snapshotHash = *ref.SnapshotHash
		}

		ui.Println(ui.Header("Diff"))

		anyDiff := false

		for _, relPath := range paths {
			absPath := filepath.Join(root, filepath.FromSlash(relPath))
			newContent, err := os.ReadFile(absPath)
			if err != nil {
				if os.IsNotExist(err) {
					newContent = nil
				} else {
					ui.Println(ui.Error(fmt.Sprintf("Failed to read %s", relPath)))
					return err
				}
			}

			var oldContent []byte
			if snapshotHash != "" {
				sf, err := coredb.GetSnapshotFile(db, snapshotHash, relPath)
				if err != nil {
					if err != coreerrors.ErrDataNotFound {
						ui.Println(ui.Error("Failed to read snapshot file"))
						return err
					}
				} else {
					obj, err := coredb.GetObject(db, sf.ObjectHash)
					if err != nil {
						ui.Println(ui.Error("Failed to read object"))
						return err
					}
					oldContent = obj.Content
				}
			}

			if bytes.Equal(oldContent, newContent) {
				continue
			}

			anyDiff = true

			if !utf8.Valid(oldContent) || !utf8.Valid(newContent) {
				fmt.Printf("diff --kuro %s\n", relPath)
				fmt.Printf("Binary files a/%s and b/%s differ\n\n", relPath, relPath)
				continue
			}

			diffText := renderSimpleDiff(relPath, oldContent, newContent)
			fmt.Print(diffText)
		}

		if !anyDiff {
			ui.Println(ui.Simple("No changes"))
		}

		return nil
	},
}

func resolveDiffPath(root, input string) (string, error) {
	absPath, err := filepath.Abs(input)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(root, absPath)
	if err != nil {
		return "", err
	}

	rel = filepath.ToSlash(rel)
	if rel == ".." || strings.HasPrefix(rel, "../") {
		return "", fmt.Errorf("path outside repository")
	}
	if rel == "." {
		return "", fmt.Errorf("path is repository root")
	}

	return rel, nil
}

func renderSimpleDiff(path string, oldContent, newContent []byte) string {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(normalizeDiffInput(oldContent)),
		B:        difflib.SplitLines(normalizeDiffInput(newContent)),
		FromFile: "a/" + path,
		ToFile:   "b/" + path,
		Context:  3,
	}

	text, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return fmt.Sprintf("diff --kuro %s\n--- a/%s\n+++ b/%s\n\n", path, path, path)
	}

	if text != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}

	return text + "\n"
}

func normalizeDiffInput(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	return s
}

func init() {
	diffCommand.Flags().StringP("file", "f", "", "diff a specific staged file")
	rootCommand.AddCommand(diffCommand)
}
