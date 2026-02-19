package cmd

import (
	"fmt"

	"github.com/greedypanda0/kuro/cli/internal/config"
	"github.com/greedypanda0/kuro/cli/internal/ui"

	coredb "github.com/greedypanda0/kuro/core/db"

	"github.com/spf13/cobra"
)

var statusCommand = &cobra.Command{
	Use:           "status",
	Short:         "Show the status",
	Long:          "Show the status of various things",
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		stageFlag, _ := cmd.Flags().GetBool("stage")

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
	statusCommand.Flags().BoolP("stage", "s", false, "show current stage files")
	rootCommand.AddCommand(statusCommand)
}
