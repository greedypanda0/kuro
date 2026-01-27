package cmd

// TODO: make workspace checkout transactional (temp dir + swap)

import (
	"fmt"

	"cli/internal/config"
	"cli/internal/repo"
	"cli/internal/ui"

	coredb "core/db"
	coreerrors "core/errors"

	"github.com/spf13/cobra"
)

var checkoutCommand = &cobra.Command{
	Use:          "checkout [branch|commit]",
	Short:        "Switch branches or restore workspace",
	Long:         "Switch branches or restore workspace to a commit snapshot",
	Args:         cobra.MaximumNArgs(1),
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

		wsFlag, _ := cmd.Flags().GetBool("ws")

		head, err := coredb.GetConfig(db, "head")
		if err != nil {
			ui.Println(ui.Error("Failed to read HEAD"))
			return err
		}

		targetBranch := head
		var snapshotHash *string
		forceWorkspace := false

		if len(args) == 0 {
			ref, err := coredb.GetRef(db, head)
			if err != nil {
				ui.Println(ui.Error("Failed to resolve HEAD"))
				return err
			}
			snapshotHash = ref.SnapshotHash
		} else {
			input := args[0]

			ref, err := coredb.GetRef(db, input)
			if err == nil {
				targetBranch = input
				snapshotHash = ref.SnapshotHash

				if targetBranch != head {
					if err := coredb.SetConfig(db, "head", targetBranch); err != nil {
						ui.Println(ui.Error("Failed to update HEAD"))
						return err
					}
				}
			} else if err == coreerrors.ErrRefNotFound {
				snapshot, err := coredb.GetSnapshot(db, input)
				if err == coreerrors.ErrSnapshotNotFound {
					ui.Println(ui.Error("Branch or commit not found"))
					return err
				}
				if err != nil {
					ui.Println(ui.Error("Failed to resolve commit"))
					return err
				}

				snapshotHash = &snapshot.Hash
				forceWorkspace = true
			} else {
				ui.Println(ui.Error("Failed to resolve branch"))
				return err
			}
		}

		if wsFlag || forceWorkspace {
			if snapshotHash == nil {
				ui.Println(ui.Simple("No commits yet"))
				return nil
			}

			if err := repo.ResetWorkspace(root, db, *snapshotHash); err != nil {
				ui.Println(ui.Error("Failed to reset workspace"))
				return err
			}

			ui.Println(ui.Success("Workspace updated"))
			return nil
		}

		if len(args) == 0 {
			ui.Println(ui.Step(fmt.Sprintf("On branch %s", head)))
			return nil
		}

		ui.Println(ui.Success("Switched to " + targetBranch))
		return nil
	},
}

func init() {
	checkoutCommand.Flags().Bool("ws", false, "reset workspace to the target snapshot")
	rootCommand.AddCommand(checkoutCommand)
}
