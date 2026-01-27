package cmd

import (
	"cli/internal/config"
	"cli/internal/ui"
	"context"
	"core/db"
	"errors"
	"strings"

	coreerrors "core/errors"

	"github.com/spf13/cobra"
)

var branchCommand = &cobra.Command{
	Use:   "branch",
	Short: "Manage branches",
	Long:  "Create, list, and delete branches",
}

var listBranchCommand = &cobra.Command{
	Use:          "list",
	Short:        "List branches",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := config.RepoRoot()
		if err != nil {
			ui.Println(ui.Error("Repository not initialized"))
			return err
		}

		database, err := db.OpenDB(config.DatabasePathFor(root))
		if err != nil {
			ui.Println(ui.Error("Failed to open repository"))
			return err
		}
		defer database.Close()

		refs, err := db.ListRefs(database)
		if err != nil {
			ui.Println(ui.Error("Failed to list branches"))
			return err
		}

		head, err := db.GetConfig(database, "head")
		if err != nil {
			ui.Println(ui.Error("Failed to read HEAD"))
			return err
		}

		for _, ref := range refs {
			if ref.Name == head {
				ui.Println(ui.ArrowRight(ref.Name))
			} else {
				ui.Println(ui.Bullet(ref.Name))
			}
		}

		return nil
	},
}

var createBranchCommand = &cobra.Command{
	Use:          "create <name>",
	Short:        "Create a new branch",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if strings.ToLower(name) == "head" {
			ui.Println(ui.Error("Invalid branch name"))
			return errors.New("Invalid branch name")
		}

		root, err := config.RepoRoot()
		if err != nil {
			ui.Println(ui.Error("Repository not initialized"))
			return err
		}

		database, err := db.OpenDB(config.DatabasePathFor(root))
		if err != nil {
			ui.Println(ui.Error("Failed to open repository"))
			return err
		}
		defer database.Close()

		created := false
		err = db.WithTx(context.Background(), database, func(tx db.DBTX) error {
			_, err = db.GetRef(tx, name)
			if err == nil {
				ui.Println(ui.Error("Branch already exists"))
				return nil
			}
			if err != coreerrors.ErrRefNotFound {
				ui.Println(ui.Error("Failed to check branch"))
				return err
			}

			head, err := db.GetConfig(tx, "head")
			if err != nil {
				ui.Println(ui.Error("Failed to read HEAD"))
				return err
			}

			currentRef, err := db.GetRef(tx, head)
			if err != nil && err != coreerrors.ErrRefNotFound {
				ui.Println(ui.Error("Failed to resolve HEAD"))
				return err
			}

			var snapshotHash *string
			if err == coreerrors.ErrRefNotFound || currentRef == nil {
				snapshotHash = nil
			} else {
				snapshotHash = currentRef.SnapshotHash
			}

			if err := db.SetRef(tx, name, snapshotHash); err != nil {
				ui.Println(ui.Error("Failed to create branch"))
				return err
			}

			created = true
			return nil
		})
		if err != nil {
			return err
		}
		if !created {
			return nil
		}

		ui.Println(ui.Success("Created branch " + name))
		return nil
	},
}

var deleteBranchCommand = &cobra.Command{
	Use:          "delete <name>",
	Short:        "Delete a branch",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if strings.ToLower(name) == "head" {
			ui.Println(ui.Error("Cannot delete HEAD"))
			return errors.New("invalid branch name")
		}

		root, err := config.RepoRoot()
		if err != nil {
			ui.Println(ui.Error("Repository not initialized"))
			return err
		}

		database, err := db.OpenDB(config.DatabasePathFor(root))
		if err != nil {
			ui.Println(ui.Error("Failed to open repository"))
			return err
		}
		defer database.Close()

		deleted := false
		err = db.WithTx(context.Background(), database, func(tx db.DBTX) error {
			head, err := db.GetConfig(tx, "head")
			if err != nil {
				ui.Println(ui.Error("Failed to read HEAD"))
				return err
			}

			if name == head {
				ui.Println(ui.Error("Cannot delete the current branch"))
				return nil
			}

			_, err = db.GetRef(tx, name)
			if err == coreerrors.ErrRefNotFound {
				ui.Println(ui.Error("Branch does not exist"))
				return nil
			}
			if err != nil {
				ui.Println(ui.Error("Failed to resolve branch"))
				return err
			}

			if err := db.DeleteRef(tx, name); err != nil {
				ui.Println(ui.Error("Failed to delete branch"))
				return err
			}

			deleted = true
			return nil
		})
		if err != nil {
			return err
		}
		if !deleted {
			return nil
		}

		ui.Println(ui.Success("Deleted branch " + name))
		return nil
	},
}

func init() {
	rootCommand.AddCommand(branchCommand)
	branchCommand.AddCommand(listBranchCommand)
	branchCommand.AddCommand(createBranchCommand)
	branchCommand.AddCommand(deleteBranchCommand)
}
