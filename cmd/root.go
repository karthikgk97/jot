package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "jot",
	Short: "A simple CLI to jot down your thoughts",
	Long: `
- Quickly noting down content:
    jot "test note"
- Using pipe to add more content:
    cat file.txt | jot "need to look at this file later"
- For showcasing default note contents:
    jot show
- You can also pass a filename path for adding to that file
    jot -f "/root/custom_directory/notes/custom_note.txt
- The same for "show":
    jot show -f "/root/custom_directory/notes/custom_note.txt
    `,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := runCommand(cmd, args)

		if err != nil {
			fmt.Println(err)
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("table", "T", "", "The table to write to. Defaults to the one provided in config file.")
	rootCmd.PersistentFlags().StringP("label", "l", "", "The label for the notes. [Required]")
	rootCmd.MarkFlagRequired("label")
	rootCmd.PersistentFlags().BoolP("high-severity", "s", false, "Boolean for High Severity. Defaults to false")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show logs")

	rootCmd.AddCommand(showCmd)
	// Set the config file name and file path
	viper.SetConfigName("config")            // Specify the config file name without extension
	viper.SetConfigType("yaml")              // Set the type of config file
	viper.AddConfigPath("$HOME/.config/jot") // Add the root directory to search for the config file

	// Read the config file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s \n", err)
		os.Exit(1)
	}
}

type JotNote struct {
	Label        string
	Content      string
	CreatedAt    time.Time
	HighSeverity bool
}

func createDBIfNotExist(dbPath string) error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		f, err := os.Create(dbPath)
		defer f.Close()

		if err != nil {
			return fmt.Errorf("DB Creation failed due to err: %v", err)
		}

	}
	return nil
}

func isTableExist(db *sql.DB, tableName string) (bool, error) {

	query := "SELECT name from sqlite_master WHERE type='table' AND name=?"
	var name string

	err := db.QueryRow(query, tableName).Scan(&name)

	if err != sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

/*
checkTableStatus function based on the configType does the following:
write: check if table exists, else creates one
read or clear: check if table exists, else returns an error
*/
func checkTableStatus(db *sql.DB, tableName string, configType string) (string, error) {
	var userConfig *viper.Viper
	switch configType {
	case "write":
		userConfig = viper.Sub("writeConfig")
	case "read":
		userConfig = viper.Sub("readConfig")
	case "clear":
		userConfig = viper.Sub("clearConfig")
	default:
		return "", errors.New("Invalid configType")
	}

	if tableName == "" {
		tableName = userConfig.GetString("defaultTable")
	}

	// check if file exists, else creating
	if tableExist, _ := isTableExist(db, tableName); !tableExist {
		if configType == "write" {
			// creating the table
			query := `
CREATE TABLE IF NOT EXISTS {{tableName}} (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    Label TEXT,
    Content TEXT,
    CreatedAt TIMESTAMP,
    HighSeverity BOOLEAN
)
`
			query = strings.Replace(query, "{{tableName}}", tableName, 1)
			_, err := db.Exec(query)

			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("Table %v does not exist. Make sure it exists before performing read or clear operations", tableName)
		}
	}

	return tableName, nil
}


func insertJotNote(db *sql.DB, jotNote JotNote, tableName string) error {
  query := `
  INSERT INTO {{tableName}} (Label, Content, CreatedAt, HighSeverity)
  VALUES (?, ?, ?, ?)
  `
  query = strings.Replace(query, "{{tableName}}", tableName, 1)

  _, err := db.Exec(query, jotNote.Label, jotNote.Content, jotNote.CreatedAt, jotNote.HighSeverity)
  
  return err
}


func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) == 0
}

// run command function
func runCommand(cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		asciiString := `
                       _    _______ 
                      | |  |__   __|
                      | | ___ | |   
                  _   | |/ _ \| |   
                 | |__| | (_) | |   
                  \____/ \___/|_|   
                                
      easy and simple cli to jot down your thoughts.
      use jot --help for more.
`
		fmt.Println("\033[33m" + asciiString + "\033[0m")
		return nil
	}

	table, _ := cmd.Flags().GetString("table")
	label, _ := cmd.Flags().GetString("label")
	sev, _ := cmd.Flags().GetBool("high-severity")

    if label == "" {
      fmt.Println("Label is empty, so using default var: 'default'")
      label = "default"
    }

	var pipeContents string
	if isInputFromPipe() {
		c, _ := io.ReadAll(os.Stdin)
		pipeContents = string(c)
	}

	contentToWrite := strings.Join(args, " ")

	if pipeContents != "" {
		contentToWrite += "\n\n"
		contentToWrite += pipeContents
	}

	jotToWrite := JotNote{
		Label:        label,
		Content:      contentToWrite,
		CreatedAt:    time.Now(),
		HighSeverity: sev,
	}

	homeDir := os.Getenv("HOME")

	if homeDir == "" {
		return fmt.Errorf("HOME Env Var not set")
	}

	dbPath := filepath.Join(homeDir, ".config/jot/jot.db")
	err := createDBIfNotExist(dbPath)

	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbPath)

    if err != nil {
      return err  
  }
    defer db.Close()

  tableName, tableErr := checkTableStatus(db, table, "write")

	if tableErr != nil {
		return tableErr
	}

  insertErr := insertJotNote(db, jotToWrite, tableName)

  if insertErr != nil {
    return insertErr
  }

    return nil
}
