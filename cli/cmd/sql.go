package cmd

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/greedypanda0/kuro/cli/internal/config"
	"github.com/greedypanda0/kuro/cli/internal/ui"
	coredb "github.com/greedypanda0/kuro/core/db"
	"github.com/spf13/cobra"
)

var sqlCommand = &cobra.Command{
	Use:          "sql <query>",
	Short:        "Run raw SQL against the repository database",
	Long:         "Run raw SQL against the repository database and print results as a table",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.TrimSpace(strings.Join(args, " "))
		if query == "" {
			ui.Println(ui.Error("Query cannot be empty"))
			return fmt.Errorf("empty query")
		}

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

		rows, err := db.Query(query)
		if err != nil {
			result, execErr := db.Exec(query)
			if execErr != nil {
				ui.Println(ui.Error("Failed to execute query"))
				return execErr
			}

			affected, _ := result.RowsAffected()
			ui.Println(ui.Success(fmt.Sprintf("Query OK. Rows affected: %d", affected)))
			return nil
		}
		defer rows.Close()

		table, err := renderTable(rows)
		if err != nil {
			ui.Println(ui.Error("Failed to render results"))
			return err
		}

		if table == "" {
			ui.Println(ui.Simple("No rows returned"))
			return nil
		}

		ui.Println(table)
		return nil
	},
}

func renderTable(rows *sql.Rows) (string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	if len(cols) > 0 {
		fmt.Fprintln(w, strings.Join(cols, "\t"))
	}

	rowCount := 0
	for rows.Next() {
		values := make([]any, len(cols))
		pointers := make([]any, len(cols))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return "", err
		}

		cells := make([]string, len(cols))
		for i, v := range values {
			switch val := v.(type) {
			case nil:
				cells[i] = "NULL"
			case []byte:
				cells[i] = string(val)
			default:
				cells[i] = fmt.Sprint(val)
			}
		}

		fmt.Fprintln(w, strings.Join(cells, "\t"))
		rowCount++
	}

	if err := rows.Err(); err != nil {
		return "", err
	}

	if err := w.Flush(); err != nil {
		return "", err
	}

	if rowCount == 0 {
		return "", nil
	}

	return buf.String(), nil
}

func init() {
	rootCommand.AddCommand(sqlCommand)
}