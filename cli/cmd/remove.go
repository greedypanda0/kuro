package cmd

import (
	"fmt"
	"path/filepath"

	"cli/internal/config"
	"cli/internal/ui"

	coredb "core/db"
	"core/ops"

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

		db, err := coredb.OpenDB(config.DatabasePath)
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

		base := filepath.Clean(path)
		stageFiles, err := coredb.GetStageFiles(db)
		if err != nil {
			ui.Println(ui.Error("Failed to get staged files"))
			return err
		}

		removed := 0

		for _, file := range stageFiles {
			if ops.IsIgnored(file.Path, []string{base}) {
				if err := coredb.RemoveStageFile(db, file.Path); err != nil {
					ui.Println(ui.Error(fmt.Sprintf("Failed to remove %s", file.Path)))
					return err
				}
				ui.Println(ui.Success(fmt.Sprintf("Removed %s from stage", file.Path)))
				removed++
			}
		}

		if removed == 0 {
			ui.Println(ui.Error(fmt.Sprintf("Path %s is not staged", path)))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCommand)
}
