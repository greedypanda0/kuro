package cmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"cli/internal/config"
	"cli/internal/ui"

	coredb "core/db"
	"core/ops"

	"github.com/spf13/cobra"
)

type objectFile struct {
	Content []byte
	Hash    string
	Path    string
}

var commitCommand = &cobra.Command{
	Use:          "commit",
	Short:        "Create a commit",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		message, _ := cmd.Flags().GetString("message")
		if strings.TrimSpace(message) == "" {
			ui.Println(ui.Error("Commit message required"))
			return errors.New("commit message required")
		}

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

		done := false

		err = coredb.WithTx(context.Background(), db, func(tx coredb.DBTX) error {
			head, err := coredb.GetConfig(tx, "head")
			if err != nil {
				ui.Println(ui.Error("Failed to get head"))
				return err
			}
			ref, err := coredb.GetRef(tx, head)
			if err != nil {
				ui.Println(ui.Error("Failed to get ref"))
				return err
			}
			stageFiles, err := coredb.GetStageFiles(tx)
			if err != nil {
				ui.Println(ui.Error("Failed to get stage files"))
				return err
			}
			if len(stageFiles) == 0 {
				ui.Println(ui.Error("No files staged"))
				return nil
			}

			objectFiles := []objectFile{}

			for _, file := range stageFiles {
				absPath := filepath.Clean(filepath.Join(root, filepath.FromSlash(file.Path)))

				content, err := os.ReadFile(absPath)
				if err != nil {
					ui.Println(ui.Error("Failed to read file"))
					return err
				}

				objectHash := ops.Hash(content)

				if err := coredb.CreateObject(tx, objectHash, content); err != nil {
					ui.Println(ui.Error("Failed to create object"))
					return err
				}

				objectFiles = append(objectFiles, objectFile{
					Path:    file.Path,
					Hash:    objectHash,
					Content: content,
				})
			}

			currentSnapshotFiles := []coredb.SnapshotFile{}

			if ref != nil && ref.SnapshotHash != nil {
				var err error
				currentSnapshotFiles, err = coredb.ListSnapshotFiles(tx, *ref.SnapshotHash)
				if err != nil {
					ui.Println(ui.Error("Failed to list snapshot files"))
					return err
				}
			}

			newSnapshotFiles := []coredb.SnapshotFile{}
			for _, file := range objectFiles {
				newSnapshotFiles = append(newSnapshotFiles, coredb.SnapshotFile{
					Path:       file.Path,
					ObjectHash: file.Hash,
				})
			}

			if coredb.CompareSnapshotFiles(currentSnapshotFiles, newSnapshotFiles) {
				ui.Println(ui.Error("No changes detected"))
				return nil
			}

			user := "root"
			var parentHash *string
			if ref != nil && ref.SnapshotHash != nil {
				parentHash = ref.SnapshotHash
			}

			sort.Slice(objectFiles, func(i, j int) bool {
				return objectFiles[i].Path < objectFiles[j].Path
			})

			var builder strings.Builder
			if parentHash != nil {
				builder.WriteString("parent:")
				builder.WriteString(*parentHash)
				builder.WriteString("\n")
			}
			builder.WriteString("message:")
			builder.WriteString(message)

			for _, obj := range objectFiles {
				builder.WriteString("\npath:")
				builder.WriteString(obj.Path)
				builder.WriteString("\nobject:")
				builder.WriteString(obj.Hash)
			}

			snapshotHash := ops.Hash([]byte(builder.String()))

			if err := coredb.CreateSnapshot(tx, snapshotHash, parentHash, message, &user); err != nil {
				ui.Println(ui.Error("Failed to create snapshot"))
				return err
			}

			for _, object := range objectFiles {
				if err := coredb.CreateSnapshotFile(tx, snapshotHash, object.Path, object.Hash); err != nil {
					ui.Println(ui.Error("Failed to create snapshot files"))
					return err
				}
			}

			if err := coredb.UpdateRef(tx, head, &snapshotHash); err != nil {
				ui.Println(ui.Error("Failed to update head ref"))
				return err
			}

			if err := coredb.ClearStage(tx); err != nil {
				ui.Println(ui.Error("Failed to clear stage"))
				return err
			}

			done = true
			return nil
		})

		if err != nil {
			return err
		}
		if !done {
			return nil
		}

		ui.Println(ui.Success("Successfully committed your changes..."))
		return nil
	},
}

func init() {
	commitCommand.Flags().StringP("message", "m", "", "commit message")
	rootCommand.AddCommand(commitCommand)
}
