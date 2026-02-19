package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/greedypanda0/kuro/cli/internal/config"
	"github.com/greedypanda0/kuro/cli/internal/ui"

	coredb "github.com/greedypanda0/kuro/core/db"
	coreerrors "github.com/greedypanda0/kuro/core/errors"

	"github.com/spf13/cobra"
)

var logsCommand = &cobra.Command{
	Use:          "logs",
	Short:        "Show commit logs",
	Long:         "Show commit logs from newest to oldest",
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

		branch, _ := cmd.Flags().GetString("branch")

		head, err := coredb.GetConfig(db, "head")
		if err != nil {
			ui.Println(ui.Error("Failed to read HEAD"))
			return err
		}

		target := head
		if branch != "" {
			target = branch
		}

		ref, err := coredb.GetRef(db, target)
		if err == coreerrors.ErrRefNotFound {
			ui.Println(ui.Error("Branch not found"))
			return err
		}
		if err != nil {
			ui.Println(ui.Error("Failed to resolve branch"))
			return err
		}

		if ref.SnapshotHash == nil {
			ui.Println(ui.Simple("No commits yet"))
			return nil
		}

		ui.Println(ui.Header("Commits"))
		ui.Println(ui.Header(fmt.Sprintf("Branch %s", ref.Name)))

		var snapshots []coredb.Snapshot
		current := *ref.SnapshotHash
		for {
			snapshot, err := coredb.GetSnapshot(db, current)
			if err == coreerrors.ErrSnapshotNotFound {
				ui.Println(ui.Warn("Commit history is incomplete"))
				return nil
			}
			if err != nil {
				ui.Println(ui.Error("Failed to read snapshot"))
				return err
			}

			snapshots = append(snapshots, *snapshot)

			if snapshot.ParentHash == nil {
				break
			}
			current = *snapshot.ParentHash
		}

		sort.Slice(snapshots, func(i, j int) bool {
			return snapshots[i].Timestamp < snapshots[j].Timestamp
		})

		for _, snapshot := range snapshots {
			timestamp := time.Unix(snapshot.Timestamp, 0).
				Format("Mon Jan 2 15:04:05 2006")
			ui.Println(ui.Step(fmt.Sprintf("%s  %s", snapshot.Hash, snapshot.Message)))
			ui.Println(ui.Muted.Render(fmt.Sprintf("  %s", timestamp)))
		}

		return nil
	},
}

func init() {
	logsCommand.Flags().StringP("branch", "b", "", "show logs for a branch")
	rootCommand.AddCommand(logsCommand)
}
