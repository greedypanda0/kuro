package cmd

import (
	"cli/internal/ui"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a kuro repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Println(ui.Step("Initializing your kuro vault...."))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
