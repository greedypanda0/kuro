package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"cli/internal/config"
	"cli/internal/ui"

	coredb "core/db"
	coreerrors "core/errors"
	"core/ops"

	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:          "add <path>",
	Short:        "Add files to the stage",
	Long:         "Add files or directories to the staging area",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		arg := args[0]

		db, err := coredb.OpenDB(config.DatabasePath)
		if err != nil {
			ui.Println(ui.Error("Failed to open database"))
			return err
		}
		defer db.Close()

		kuroIgnore, err := ops.ReadKuroIgnore(config.IgnorePath)
		if err != nil {
			if err == coreerrors.ErrIgnoreFileNotFound {
				kuroIgnore = []string{}
			} else {
				ui.Println(ui.Error("Failed to read .kuroignore"))
				return err
			}
		}

		var path string
		if arg == "." {
			cwd, err := os.Getwd()
			if err != nil {
				ui.Println(ui.Error("Failed to get current directory"))
				return err
			}
			path = cwd
		} else {
			path = arg
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			ui.Println(ui.Error("Failed to resolve path"))
			return err
		}

		info, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				ui.Println(ui.Error(fmt.Sprintf("Path does not exist: %s", arg)))
				return err
			}
			ui.Println(ui.Error("Failed to stat path"))
			return err
		}

		filesToStage := []string{}

		if info.IsDir() {
			ui.Println(ui.Step("Scanning directory..."))

			files, err := ops.ReadDir(absPath)
			if err != nil {
				ui.Println(ui.Error("Failed to read directory"))
				return err
			}

			for _, file := range files {
				if ops.IsIgnored(file.Path, kuroIgnore) {
					continue
				}
				filesToStage = append(filesToStage, file.Path)
			}
		} else {
			if !ops.IsIgnored(absPath, kuroIgnore) {
				filesToStage = append(filesToStage, absPath)
			}
		}

		total := len(filesToStage)
		if total == 0 {
			ui.Println(ui.Step("Nothing to stage"))
			return nil
		}

		ui.Println(ui.Step(fmt.Sprintf("Staging %d file(s)...", total)))

		for i, file := range filesToStage {
			ratio := float64(i+1) / float64(total)
			ui.Println(ui.Progress(30, ratio))

			if err := coredb.AddStageFile(db, file); err != nil {
				ui.Println(ui.Error(fmt.Sprintf("Failed to stage file: %s", file)))
				return err
			}
		}

		ui.Println(ui.Success(fmt.Sprintf("Staged %d file(s)", total)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCommand)
}
