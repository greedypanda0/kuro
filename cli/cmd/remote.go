package cmd

import (
	"errors"
	"strings"

	"github.com/greedypanda0/kuro/cli/internal/config"
	"github.com/greedypanda0/kuro/cli/internal/ui"
	"github.com/greedypanda0/kuro/core/db"
	coreerrors "github.com/greedypanda0/kuro/core/errors"
	"github.com/spf13/cobra"
)

var remoteCommand = &cobra.Command{
	Use:   "remote",
	Short: "Manage remote",
	Long:  "Manage remote",
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

		remote, err := db.GetConfig(database, "remote")
		if errors.Is(err, coreerrors.ErrDataNotFound) {
			ui.Println(ui.Error("Remote not found"))
			return nil
		} else if err != nil {
			ui.Println(ui.Error("Failed to get remote"))
			return err
		}

		ui.Println(ui.ArrowRight(remote))
		return nil
	},
}

var remoteAddCommand = &cobra.Command{
	Use:   "add <name>",
	Short: "Add remote",
	Long:  "Add remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			ui.Println(ui.Error("Invalid arguments"))
			return errors.New("invalid arguments")
		}
		remote := args[0]

		root, err := config.RepoRoot()
		if err != nil {
			ui.Println(ui.Error("Repository not initialized"))
			return err
		}

		trimmed := strings.Trim(remote, "/")
		parts := strings.Split(trimmed, "/")
		if len(parts) < 2 {
			ui.Println(ui.Error("Invalid remote format. Use <user>/<repo> or a URL ending with that."))
			return errors.New("invalid remote format")
		}

		repo := parts[len(parts)-1]
		user := parts[len(parts)-2]

		database, err := db.OpenDB(config.DatabasePathFor(root))
		if err != nil {
			ui.Println(ui.Error("Failed to open repository"))
			return err
		}
		defer database.Close()

		_, configError := db.GetConfig(database, "remote")
		if configError == nil {
			ui.Println(ui.Error("Remote already exists"))
			return errors.New("remote already exists")
		}
		if !errors.Is(configError, coreerrors.ErrDataNotFound) {
			ui.Println(ui.Error("Failed to check remote"))
			return configError
		}

		if err := db.SetConfig(database, "remote", user+"/"+repo); err != nil {
			ui.Println(ui.Error("Failed to add remote"))
			return err
		}

		ui.Println(ui.Success("Added remote "), ui.ArrowRight(remote))
		return nil
	},
}

var remoteRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove remote",
	Long:  "Remove remote",
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

		if err := db.DeleteConfig(database, "remote"); err != nil {
			ui.Println(ui.Error("Failed to remove remote"))
			return err
		}

		ui.Println(ui.Success("Removed remote!"))
		return nil
	},
}

func init() {
	remoteCommand.AddCommand(remoteAddCommand)
	remoteCommand.AddCommand(remoteRemoveCommand)
	rootCommand.AddCommand(remoteCommand)
}
