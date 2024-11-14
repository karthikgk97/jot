package cmd

import (
	// "errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	rootCmd.PersistentFlags().StringP("table", "t", "", "The table to write to. Defaults to the one provided in config file.")
	rootCmd.PersistentFlags().BoolP("high", "H", false, "Boolean for High Severity. Defaults to false. There is no --low flag as default is low severity")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show logs")

	rootCmd.AddCommand(uiCmd)
	// rootCmd.AddCommand(clearCmd)
	// rootCmd.AddCommand(statCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
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

// inserts note to a file
func insertJotNote(folderPath string, fileContents string) error {

	currentTime := time.Now().Format("2006_01_02__15_04_05")

	fileToWrite := filepath.Join(folderPath, currentTime+".txt")
	err := os.WriteFile(fileToWrite, []byte(fileContents), 0644)

	if err != nil {
		return fmt.Errorf("Errored while writing to file with error: %s", err)
	}

	return nil
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) == 0
}

// run command function
func runCommand(cmd *cobra.Command, args []string) error {

	if !isInputFromPipe() && len(args) == 0 {
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
		fmt.Println("\033[36m" + asciiString + "\033[0m")
		return nil
	}

	highSeverity, _ := cmd.Flags().GetBool("high")
	verbose, _ := cmd.Flags().GetBool("verbose")

	if verbose == true {
		fmt.Println("Verbose is set to True")
		fmt.Println("High Severity: ", highSeverity)
	}
	var pipeContents string
	if isInputFromPipe() {
		fmt.Println("Grabbing contents from pipe operation output")
		c, _ := io.ReadAll(os.Stdin)
		pipeContents = string(c)
		fmt.Println("Contents from pipe operation:\n", pipeContents)
	}

	contentToWrite := strings.Join(args, " ")

	if pipeContents != "" {
		contentToWrite += "\n\n"
		contentToWrite += pipeContents
	}

	fmt.Println("Contents to Write is ", contentToWrite)

	homeDir := os.Getenv("HOME")

	if homeDir == "" {
		return fmt.Errorf("HOME Env Var not set")
	}

	FolderPath := filepath.Join(homeDir, ".config/jot/")

	insertJotNote(FolderPath, contentToWrite)

	fmt.Println("Jotted!")

	return nil
}
