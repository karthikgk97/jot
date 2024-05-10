package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show stats",
	Long: `
  Stats allows you to get metrics about your Jot Notes.
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := showStatsForJotNote()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func showStatsForJotNote() error {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return fmt.Errorf("HOME Env Var not set")
	}

	dbPath := filepath.Join(homeDir, ".config/jot/jot.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	numTablesQuery := "SELECT COUNT(*) AS total_tables FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';"

	var totalTables int
	err = db.QueryRow(numTablesQuery).Scan(&totalTables)

	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Println("-------------------------------")
	fmt.Println("|\033[36mTotal number of Tables: \033[0m", totalTables, "  |")
	fmt.Println("-------------------------------")
    fmt.Println()

	tableNameQuery := "SELECT name AS table_name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"

	tables, err := db.Query(tableNameQuery)

	defer tables.Close()
	if err != nil {
		return err
	}

	for tables.Next() {
		var t string
		err := tables.Scan(&t)
		if err != nil {
			return err
		}

		fmt.Println("\033[36mFor table: \033[0m", t)

		rowCountQuery := "SELECT COUNT(*) FROM " + t

		var tableRowCount int
		tableRowErr := db.QueryRow(rowCountQuery).Scan(&tableRowCount)

		if tableRowErr != nil {
			return tableRowErr
		}

		fmt.Println("\033[38;5;208m    Total # of Rows: \033[0m", tableRowCount)

		statsQuery := "SELECT Label, HighSeverity, Count(*) FROM " + t + " GROUP BY Label, HighSeverity"
		labelRows, err := db.Query(statsQuery)
		if err != nil {
			return err
		}

		defer labelRows.Close()

		for labelRows.Next() {
			var label string
			var highSeverity bool
			var count int

			err := labelRows.Scan(&label, &highSeverity, &count)
			if err != nil {
				return err
			}

			sev := "\033[32mlow"

			if highSeverity == true {
				sev = "\033[31mhigh"
			}
			fmt.Printf("\033[36m    Label: \033[0m%v \033[36m| Severity: \033[0m%v \033[36m| Count: \033[0m %v\n", label, sev, count)

		}
		fmt.Println("\033[33m----------------------------------------------------------------------\033[0m")

	}
	return nil
}
