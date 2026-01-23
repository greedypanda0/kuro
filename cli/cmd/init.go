package cmd

import (
	stderrors "errors"

	"cli/internal/config"
	"cli/internal/ui"
	"core/db"
	kuroerrors "core/errors"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a kuro repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Println(ui.Step("Initializing your kuro repository..."))

		database, err := db.InitSQL(config.DatabasePath)
		if err != nil {
			if stderrors.Is(err, kuroerrors.ErrRepoAlreadyInitialized) {
				ui.Println(ui.Error("Repository already exists"))
				return nil
			}

			ui.Println(ui.Error("Failed to initialize kuro repository"))
			return err
		}
		defer database.Close()

		if err := db.ApplySchema(database); err != nil {
			ui.Println(ui.Error("Failed to apply schema"))
			return err
		}

		ui.Println(ui.Success("Created your kuro repository!"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
