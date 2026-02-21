package cmd

import (
	"github.com/greedypanda0/kuro/cli/internal/config"
	"github.com/greedypanda0/kuro/cli/internal/ui"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var whoamicommand = &cobra.Command{
	Use:          "whoami",
	Short:        "check your auth token",
	Long:         "check who are you",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			ui.Println(ui.Error("Failed to load config."))
			return err
		}

		if cfg.Name != "" {
			ui.Println(ui.Bullet("You are " + cfg.Name))
		} else {
			ui.Println(ui.Cross("You have no name, add one by config"))

		}

		if cfg.Token != "" {
			req, err := http.NewRequest("GET", config.ApiUrl+"/users/me", nil)
			if err != nil {
				ui.Println(ui.Error("Failed to create request."))
				return err
			}

			req.Header.Set("Authorization", "Bearer "+cfg.Token)
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				ui.Println(ui.Error("Failed to send request."))
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				ui.Println(ui.Error(fmt.Sprintf("request failed: %s", resp.Status)))
				return err
			}

			var data struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
				return err
			}

			ui.Println(ui.Bullet(fmt.Sprintf("You are %s {%s} [remote]", data.Name, data.Email)))
		} else {
			ui.Println(ui.Cross("You have no token added, add one by config"))
		}

		return nil
	},
}

func init() {
	rootCommand.AddCommand(whoamicommand)
}
