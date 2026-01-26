package cmd

import (
	"fmt"

	"cli/internal/config"
	"cli/internal/ui"

	coredb "core/db"

	"github.com/spf13/cobra"
)

var StatusCommand = &cobra.Command{
	Use:           "status",
	Short:         "Show the status",
	Long:          "Show the status of various things",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		stageFlag, _ := cmd.Flags().GetBool("stage")

		db, err := coredb.OpenDB(config.DatabasePath)
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

		ui.Println(ui.Step(fmt.Sprintf("On branch %s", head)))

		if ref.SnapshotHash == nil {
			ui.Println(ui.Simple("No commits yet"))
		} else {
			ui.Println(ui.Step(fmt.Sprintf("Commit: %s", *ref.SnapshotHash)))
		}

		if stageFlag {
			stageFiles, err := coredb.GetStageFiles(db)
			if err != nil {
				ui.Println(ui.Error("Failed to get staged files"))
				return err
			}

			ui.Println(ui.Header("Staged files"))

			if len(stageFiles) == 0 {
				ui.Println(ui.Simple("No files staged"))
				return nil
			}

			for _, file := range stageFiles {
				ui.Println(ui.Bullet(file.Path))
			}

			ui.Println(ui.Step(fmt.Sprintf("Total: %d", len(stageFiles))))
		}

		return nil
	},
}

func init() {
	StatusCommand.Flags().BoolP("stage", "s", false, "show current stage files")
	rootCmd.AddCommand(StatusCommand)
}
