package cmd

import (
	"bytes"
	"cli/internal/config"
	"cli/internal/ui"
	coredb "core/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

		db, err := coredb.OpenDB(config.DatabasePathFor(root))
		if err != nil {
			ui.Println(ui.Error("Failed to open repository"))
			return err
		}
		defer db.Close()

		cfg, err := config.LoadConfig()
		if err != nil {
			ui.Println(ui.Error("Failed to load config"))
			return err
		}
		if cfg.Token == "" {
			ui.Println(ui.Error("No token found"))
			return nil
		}
		wd, err := os.Getwd()
		if err != nil {
			ui.Println(ui.Error("Failed to get current working directory"))
			return err
		}

		rootName := filepath.Base(wd)

		head, err := coredb.GetConfig(db, "head")
		if err != nil {
			ui.Println(ui.Error("Failed to read head ref"))
			return err
		}
		ref, err := coredb.GetRef(db, head)
		if err != nil {
			ui.Println(ui.Error("Failed to resolve head ref"))
			return err
		}

		// snapshot
		snapshot, err := coredb.GetSnapshot(db, *ref.SnapshotHash)
		if err != nil {
			ui.Println(ui.Error("Failed to get snapshot"))
			return err
		}

		// files
		files, err := coredb.ListSnapshotFiles(db, snapshot.Hash)
		if err != nil {
			ui.Println(ui.Error("Failed to list snapshot files"))
			return err
		}

		// objects
		var objects []*coredb.Object
		for _, file := range files {
			object, err := coredb.GetObject(db, file.ObjectHash)
			if err != nil {
				ui.Println(ui.Error("Failed to get object"))
				return err
			}
			objects = append(objects, object)
		}

		payload := map[string]any{
			"metadata": map[string]string{
				"name": rootName,
				"head": head,
			},
			"ref":      ref,
			"snapshot": snapshot,
			"files":    files,
			"objects":  objects,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			ui.Println(ui.Error("Failed to encode request body"))
			return err
		}
		body := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", config.ApiUrl+"/repositories", body)
		if err != nil {
			ui.Println(ui.Error("Failed to create request"))
			return err
		}
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ui.Println(ui.Error("Failed to send request"))
			return err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			ui.Println(ui.Error("Failed to read response"))
			return err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			ui.Println(ui.Error("Push failed: " + string(respBody)))
			return fmt.Errorf("push failed with status %s", resp.Status)
		}

		if len(respBody) > 0 {
			ui.Println(string(respBody))
		}

		return nil
	},
}

func init() {
	rootCommand.AddCommand(pushCommand)
}
