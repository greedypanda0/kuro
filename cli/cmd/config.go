package cmd

import (
	"github.com/greedypanda0/kuro/cli/internal/config"
	"github.com/greedypanda0/kuro/cli/internal/ui"
	"os"

	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:          "config",
	Short:        "set config",
	Long:         "set config for your env",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		token, _ := cmd.Flags().GetString("token")
		
		cfg, err := config.LoadConfig()
		if err != nil {
			if os.IsNotExist(err) {
				cfg = &config.Config{}
			} else {
				ui.Println(ui.Error("failed to load config"))
				return err
			}
		}

		if name != "" {
			cfg.Name = name
		}
		if token != "" {
			cfg.Token = token
		}
		
		if err := config.SaveConfig(cfg); err != nil {
			ui.Println(ui.Error("Failed to save config."))
			return err
		}

		ui.Println(ui.Success("Config saved successfully."))
		return nil
	},
}

func init() {
	configCommand.Flags().String("name", "", "set name")
	configCommand.Flags().String("token", "", "set auth token")
	rootCommand.AddCommand(configCommand)
}
