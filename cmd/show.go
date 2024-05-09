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

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show notes",
	Long: `
  Shows allows you to view notes.
  By Default: Shows N notes (oldest or recent notes is based off of decent config).
  N is either the one provided in argument or default var in config`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := showJotNote()
		if err != nil {
			fmt.Println(err)
		}
	},
}

var (
	table          string
	label          string
	numNotes       int64
	afterTimeline  string
	beforeTimeline string
	high           bool
	low            bool
	showOldest     bool
	showRecent     bool
)

func init() {
	showCmd.PersistentFlags().StringVarP(&table, "table", "t", "", "The table to show data from. Defaults to the one provided in config file.")
	showCmd.PersistentFlags().StringVarP(&label, "label", "l", "", `
    The label for the notes. Defaults to the one provided in config file.
    For no label filter, pass label='no-label' (Assuming no-label is also not a label passed before)
`)
	showCmd.PersistentFlags().Int64VarP(&numNotes, "num-notes", "n", 0, "The number of notes to display. Defaults to the one provided in config file")

	showCmd.PersistentFlags().StringVar(&afterTimeline, "after", "", "The timeline after to get the notes for. In YYYY-MM-DD format")
	showCmd.PersistentFlags().StringVar(&beforeTimeline, "before", "", "The timeline before to get the notes for. In YYYY-MM-DD format")
	showCmd.PersistentFlags().BoolVar(&high, "high", false, "Boolean for High Only Severity. Defaults to false")
	showCmd.PersistentFlags().BoolVar(&low, "low", false, "Boolean for Low Only Severity. Defaults to false")

	showCmd.PersistentFlags().BoolVar(&showOldest, "oldest", false, "Boolean for Showing Oldest N notes. Defaults to false")
	showCmd.PersistentFlags().BoolVar(&showRecent, "recent", false, "Boolean for Showing Recent N notes. Defaults to false")
}

func showQueryBuilder(tableName string, viewPreference string) string {

	var queryFilter []string

	if label == "" {
		label = viper.GetString("showConfig.defaultLabel")
	}

	if label != "no-label" {
		queryFilter = append(queryFilter, fmt.Sprintf("Label = '%s'", label))
	}

	if afterTimeline != "" {

		queryFilter = append(queryFilter, fmt.Sprintf("CreatedAt > '%s'", afterTimeline))
	}

	if beforeTimeline != "" {
		queryFilter = append(queryFilter, fmt.Sprintf("CreatedAt < '%s'", beforeTimeline))
	}

	if high {
		queryFilter = append(queryFilter, "HighSeverity = 1")
	}

	if low {
		queryFilter = append(queryFilter, "HighSeverity = 0")
	}

	var orderBy string
	var revOrderBy string
	switch viewPreference {
	case "recent":
		orderBy = "DESC"
		revOrderBy = "ASC"
	case "oldest":
		orderBy = "ASC"
		revOrderBy = "DESC"
	default:
		orderBy = "DESC"
		revOrderBy = "ASC"
	}

	if numNotes == 0 {
		numNotes = viper.GetInt64("showConfig.notesToDisplay")
	}

	var query string
	query = fmt.Sprintf("SELECT Label, Content, HighSeverity, CreatedAt FROM %v", tableName)

	if len(queryFilter) > 0 {
		query += " WHERE " + strings.Join(queryFilter, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY CreatedAt %v LIMIT %v", orderBy, numNotes)

	modifiedQuery := "SELECT Label, Content, HighSeverity FROM  (" + query + ") AS  subquery " + "ORDER BY CreatedAt " + revOrderBy

	return modifiedQuery
}

func showJotNote() error {
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

	tableName, tableErr := EnsureTableExists(db, table, "show")
	if tableErr != nil {
		return tableErr
	}

	if showOldest && showRecent {
		return fmt.Errorf("Cant show both Oldest and Recent. Please pass one of the args")
	}

	var viewPref string
	if showOldest {
		viewPref = "oldest"
	} else if showRecent {
		viewPref = "recent"
	} else {
		viewPref = viper.GetString("showConfig.defaultViewPreference")
	}

	queryToExec := showQueryBuilder(tableName, viewPref)

	rows, err := db.Query(queryToExec)

	if err != nil {
		return err
	}

	for rows.Next() {
		var content string
		var label string
		var sev bool

		err := rows.Scan(&label, &content, &sev)
		if err != nil {
			fmt.Println(err)
			return err
		} else {

			var s string
			if sev == false {
				s = "low"
			} else {
				s = "\033[31mhigh\033[0m"
			}

			var outputString string
			outputString += "\033[36mLabel: \033[0m" + label + "\n"
			outputString += "\033[36mContent: \033[0m\n\n" + content + "\n"
			outputString += "\033[36mSeverity: \033[0m" + s + "\n"
			outputString += "\033[33m----------------------------------------------------------------------\033[0m" + "\n"
			fmt.Println(outputString)
		}
	}

	defer rows.Close()

	return nil
}
