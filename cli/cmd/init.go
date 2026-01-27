package cmd

import (
	stderrors "errors"
	"os"
	"path/filepath"

	"cli/internal/config"
	"cli/internal/ui"
	"core/db"
	kuroerrors "core/errors"

	"github.com/spf13/cobra"
)

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "Initialize a kuro repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Println(ui.Step("Initializing your kuro repository..."))

		root, err := os.Getwd()
		if err != nil {
			ui.Println(ui.Error("Failed to resolve repository root"))
			return err
		}

		if _, err := os.Stat(filepath.Join(root, config.RepoDir)); err == nil {
			ui.Println(ui.Error("Repository already exists"))
			return nil
		} else if !os.IsNotExist(err) {
			ui.Println(ui.Error("Failed to check repository"))
			return err
		}

		database, err := db.InitSQL(config.DatabasePathFor(root))
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

		if err := InitIgnore(config.IgnorePathFor(root)); err != nil {
			ui.Println(ui.Error("Failed to initialize .kuroignore"))
			return err
		}

		ui.Println(ui.Success("Created your kuro repository!"))
		return nil
	},
}

func InitIgnore(ignorePath string) error {
	path := filepath.Join(ignorePath)

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	content := []byte(".kuro\n.git\nnode_modules\ndist\nbuild\n\n")

	return os.WriteFile(path, content, 0644)
}

func init() {
	rootCommand.AddCommand(initCommand)
}
