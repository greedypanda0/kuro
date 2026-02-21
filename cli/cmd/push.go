package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/greedypanda0/kuro/cli/internal/config"
	"github.com/greedypanda0/kuro/cli/internal/ui"
	"github.com/greedypanda0/kuro/core/db"
	coreerrors "github.com/greedypanda0/kuro/core/errors"

	"github.com/spf13/cobra"
)

var pushCommand = &cobra.Command{
	Use:          "push",
	Short:        "push the changes to remote",
	Long:         "push the changes to remote server",
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

		cfg, err := config.LoadConfig()
		if err != nil {
			ui.Println(ui.Error("Failed to load config"))
			return err
		}
		if cfg.Token == "" {
			ui.Println(ui.Error("No token found"))
			return nil
		}

		remote, err := db.GetConfig(database, "remote")
		if errors.Is(err, coreerrors.ErrDataNotFound) {
			ui.Println(ui.Error("Remote not found"))
			return nil
		} else if err != nil {
			ui.Println(ui.Error("Failed to get remote"))
			return err
		}

		file, err := os.Open(config.DatabasePathFor(root))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		req, err := http.NewRequest("POST", config.ApiUrl+"/repositories", file)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Token))
		req.Header.Set("X-Remote", remote)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		ui.Println(ui.Success("Successfully pushed the .db"))

		return nil
	},
}

func init() {
	rootCommand.AddCommand(pushCommand)
}
