package cmd

import (
	"fmt"
	"sort"

	"github.com/greedypanda0/kuro/cli/internal/config"
	"github.com/greedypanda0/kuro/cli/internal/ui"

	coredb "github.com/greedypanda0/kuro/core/db"
	coreerrors "github.com/greedypanda0/kuro/core/errors"
	"github.com/greedypanda0/kuro/core/ops"

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

			stagedSet := make(map[string]struct{}, len(stageFiles))
			if len(stageFiles) == 0 {
				ui.Println(ui.Simple("No files staged"))
			} else {
				for _, file := range stageFiles {
					stagedSet[file.Path] = struct{}{}
					ui.Println(ui.Bullet(file.Path))
				}
				ui.Println(ui.Step(fmt.Sprintf("Total: %d", len(stageFiles))))
			}

			kuroIgnore, err := ops.ReadKuroIgnore(config.IgnorePathFor(root))
			if err == coreerrors.ErrIgnoreFileNotFound {
				kuroIgnore = []string{}
			} else if err != nil {
				ui.Println(ui.Error("Failed to read ignore file"))
				return err
			}

			files, err := ops.ReadDir(root)
			if err != nil {
				ui.Println(ui.Error("Failed to read workspace"))
				return err
			}

			var unstaged []string
			for _, file := range files {
				if ops.IsIgnored(file.Path, kuroIgnore) {
					continue
				}
				if _, ok := stagedSet[file.Path]; ok {
					continue
				}
				unstaged = append(unstaged, file.Path)
			}

			sort.Strings(unstaged)

			ui.Println(ui.Header("Unstaged files"))

			if len(unstaged) == 0 {
				ui.Println(ui.Simple("No unstaged files"))
			} else {
				for _, file := range unstaged {
					ui.Println(ui.Simple("- " + file))
				}
			}
		}

		return nil
	},
}

func init() {
	statusCommand.Flags().BoolP("stage", "s", false, "show stage")
	rootCommand.AddCommand(statusCommand)
}
