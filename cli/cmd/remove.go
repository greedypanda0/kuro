package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"cli/internal/config"
	"cli/internal/ui"

	coredb "core/db"

	"github.com/spf13/cobra"
)

var removeCommand = &cobra.Command{
	Use:          "remove <path>",
	Short:        "Remove a file or directory from stage",
	Long:         "Remove a file or directory from the staging area",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

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

		if path == "." {
			if err := coredb.ClearStage(db); err != nil {
				ui.Println(ui.Error("Failed to clear stage"))
				return err
			}
			ui.Println(ui.Success("Cleared the stage"))
			return nil
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

		base := filepath.ToSlash(filepath.Clean(relToRoot))
		removed := 0

		err = coredb.WithTx(context.Background(), db, func(tx coredb.DBTX) error {
			stageFiles, err := coredb.GetStageFiles(tx)
			if err != nil {
				ui.Println(ui.Error("Failed to get staged files"))
				return err
			}

			for _, file := range stageFiles {
				if file.Path == base || strings.HasPrefix(file.Path, base+"/") {
					if err := coredb.RemoveStageFile(tx, file.Path); err != nil {
						return err
					}
					ui.Println(ui.Success(fmt.Sprintf("Removed %s from stage", file.Path)))
					removed++
				}
			}
			return nil
		})
		if err != nil {
			ui.Println(ui.Error("Failed to remove staged files"))
			return err
		}

		if removed == 0 {
			ui.Println(ui.Error(fmt.Sprintf("Path %s is not staged", path)))
		}

		return nil
	},
}

func init() {
	rootCommand.AddCommand(removeCommand)
}
