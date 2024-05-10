package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear notes",
	Long: `
  clear allows you to Clear/Delete the notes.
  By default: Clears the N notes (oldest or newest is based on what is set in config).
  N is either the one provided in argument or default var in config`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := clearJotNote()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var (
	cTable          string
	cLabel          string
	cNumNotes       int64
	cAfterTimeline  string
	cBeforeTimeline string
	cHigh           bool
	cLow            bool
	cShowOldest     bool
	cShowRecent     bool
	cRowID          int
)

func init() {
	clearCmd.PersistentFlags().StringVarP(&cTable, "table", "t", "", "The table to clear data for. Defaults to the one provided in config file.")
	clearCmd.PersistentFlags().StringVarP(&cLabel, "label", "l", "", `
    The label for the notes. Defaults to the one provided in config file.
    For no label filter, pass label='no-label' (Assuming no-label is also not a label passed before)
`)
	clearCmd.PersistentFlags().Int64VarP(&cNumNotes, "num-notes", "n", 0, "The number of notes to clear. Defaults to the one provided in config file")

	clearCmd.PersistentFlags().StringVar(&cAfterTimeline, "after", "", "The timeline after to get the notes for. In YYYY-MM-DD format")
	clearCmd.PersistentFlags().StringVar(&cBeforeTimeline, "before", "", "The timeline before to get the notes for. In YYYY-MM-DD format")
	clearCmd.PersistentFlags().BoolVar(&cHigh, "high", false, "Boolean for High Only Severity. Defaults to false")
	clearCmd.PersistentFlags().BoolVar(&cLow, "low", false, "Boolean for Low Only Severity. Defaults to false")

	clearCmd.PersistentFlags().BoolVar(&cShowOldest, "oldest", false, "Boolean for Showing Oldest N notes. Defaults to false")
	clearCmd.PersistentFlags().BoolVar(&cShowRecent, "recent", false, "Boolean for Showing Recent N notes. Defaults to false")
	clearCmd.PersistentFlags().IntVar(&cRowID, "row-id", 0, "Clearing a specific ROW ID. When passed, other filter arguments are not considered")
}

func clearQueryBuilder(tableName string, viewPreference string) string {
	var queryFilter []string

    if cRowID != 0{
      query := fmt.Sprintf("DELETE FROM %v WHERE row_id = %v", tableName, cRowID)
      return query
    }

	if cLabel == "" {
		cLabel = viper.GetString("clearConfig.defaultLabel")
	}

	if cLabel != "no-label" {
		queryFilter = append(queryFilter, fmt.Sprintf("Label = '%s'", cLabel))
	}

	if cAfterTimeline != "" {
		queryFilter = append(queryFilter, fmt.Sprintf("CreatedAt > '%s'", cAfterTimeline))
	}

	if cBeforeTimeline != "" {
		queryFilter = append(queryFilter, fmt.Sprintf("CreatedAt < '%s'", cBeforeTimeline))
	}

	if cHigh {
		queryFilter = append(queryFilter, "HighSeverity = 1")
	}

	if cLow {
		queryFilter = append(queryFilter, "HighSeverity = 0")
	}

	var orderBy string
	switch viewPreference {
	case "recent":
		orderBy = "DESC"
	case "oldest":
		orderBy = "ASC"
	default:
		orderBy = "DESC"
	}

	if cNumNotes == 0 {
		cNumNotes = viper.GetInt64("clearConfig.notesToClear")
	}

	var query string

	if len(queryFilter) > 0 {
		query = " WHERE " + strings.Join(queryFilter, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY CreatedAt %v LIMIT %v", orderBy, cNumNotes)

	modifiedQuery := "DELETE FROM " + tableName + " WHERE row_id IN ( SELECT row_id FROM " + tableName + " " + query + ");"

	return modifiedQuery
}

func clearJotNote() error {

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

	tableName, tableErr := EnsureTableExists(db, cTable, "clear")
	if tableErr != nil {
		return tableErr
	}

	if cShowOldest && cShowRecent {
		return fmt.Errorf("ERROR: Cant show both Oldest and Recent. Please pass one of the args")
	}

	var viewPref string
	if cShowOldest {
		viewPref = "oldest"
	} else if cShowRecent {
		viewPref = "recent"
	} else {
		viewPref = viper.GetString("clearConfig.defaultClearPreference")
	}

	queryToExec := clearQueryBuilder(tableName, viewPref)

	_, dbErr := db.Exec(queryToExec)

	if dbErr != nil {
		return dbErr
	}

	fmt.Println("Jot Erased")

	return nil
}
