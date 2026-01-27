package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cli/internal/config"
	"cli/internal/ui"

	coredb "core/db"
	coreerrors "core/errors"
	"core/ops"

	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:          "add <path>",
	Short:        "Add files to the stage",
	Long:         "Add files or directories to the staging area",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		arg := args[0]

		root, err := config.RepoRoot()
		if err != nil {
			ui.Println(ui.Error("Repository not initialized"))
			return err
		}

		db, err := coredb.OpenDB(config.DatabasePathFor(root))
		if err != nil {
			ui.Println(ui.Error("Failed to open database"))
			return err
		}
		defer db.Close()

		kuroIgnore, err := ops.ReadKuroIgnore(config.IgnorePathFor(root))
		if err == coreerrors.ErrIgnoreFileNotFound {
			kuroIgnore = []string{}
		}

		var path string
		if arg == "." {
			path = root
		} else {
			path = arg
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			ui.Println(ui.Error("Failed to resolve path"))
			return err
		}

		relToRoot, err := filepath.Rel(root, absPath)
		if err != nil {
			ui.Println(ui.Error("Failed to resolve repository path"))
			return err
		}
		relToRoot = filepath.ToSlash(relToRoot)
		if relToRoot == ".." || strings.HasPrefix(relToRoot, "../") {
			ui.Println(ui.Error("Path is outside the repository"))
			return fmt.Errorf("path outside repository")
		}

		info, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				ui.Println(ui.Error(fmt.Sprintf("Path does not exist: %s", arg)))
				return err
			}
			ui.Println(ui.Error("Failed to stat path"))
			return err
		}

		filesToStage := []string{}

		if info.IsDir() {
			ui.Println(ui.Step("Scanning directory..."))

			files, err := ops.ReadDir(absPath)
			if err != nil {
				ui.Println(ui.Error("Failed to read directory"))
				return err
			}

			for _, file := range files {
				relPath := filepath.ToSlash(filepath.Join(relToRoot, file.Path))
				if ops.IsIgnored(relPath, kuroIgnore) {
					continue
				}
				filesToStage = append(filesToStage, relPath)
			}
		} else {
			if !ops.IsIgnored(relToRoot, kuroIgnore) {
				filesToStage = append(filesToStage, relToRoot)
			}
		}

		filtered := make([]string, 0, len(filesToStage))
		for _, file := range filesToStage {
			absPath := filepath.Clean(filepath.Join(root, filepath.FromSlash(file)))
			content, err := os.ReadFile(absPath)
			if err != nil {
				ui.Println(ui.Error("Failed to read file"))
				return err
			}

			objectHash := ops.Hash(content)
			_, err = coredb.GetObject(db, objectHash)
			if err == nil {
				continue
			}
			if !errors.Is(err, coreerrors.ErrObjectNotFound) {
				ui.Println(ui.Error("Failed to check object"))
				return err
			}

			filtered = append(filtered, file)
		}

		filesToStage = filtered
		total := len(filesToStage)
		if total == 0 {
			ui.Println(ui.Step("Nothing to stage"))
			return nil
		}

		ui.Println(ui.Step(fmt.Sprintf("Staging %d file(s)...", total)))

		err = coredb.WithTx(context.Background(), db, func(tx coredb.DBTX) error {
			for i, file := range filesToStage {
				ratio := float64(i+1) / float64(total)
				ui.Println(ui.Progress(30, ratio))

				if err := coredb.AddStageFile(tx, file); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			ui.Println(ui.Error("Failed to stage files"))
			return err
		}

		ui.Println(ui.Success(fmt.Sprintf("Staged %d file(s)", total)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCommand)
}
