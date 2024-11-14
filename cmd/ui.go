package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Interact in UI mode",
	Long: `
    UI allows you to interact with Jot in a interactive TUI mode`,
	Args: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := showUI()
		if err != nil {
			fmt.Println(err)
		}
	},
}

type JotFile struct {
	fileName     string
	fileContents string
}

func readFiles(folderPath string) ([]JotFile, error) {
	rootFile, err := os.Open(folderPath)
	if err != nil {
		fmt.Println("Errored during opening the folderPath: ", err)
		return nil, err

	}

	defer rootFile.Close()
	filesList, err := rootFile.ReadDir(-1)

	if err != nil {
		fmt.Println("Errored reading rootFile dir with error: ", err)
		return nil, err
	}

	var jotFiles []JotFile

	for _, f := range filesList {
		if filepath.Ext(f.Name()) != ".txt" {
			continue
		}
		fName := f.Name()
		filePath := filepath.Join(folderPath, f.Name())

		fileContents, err := os.ReadFile(filePath)

		if err != nil {
			fmt.Println("Failed while reading file contents with error", err)
			return nil, err
		}

		jotFiles = append(jotFiles, JotFile{fileName: fName, fileContents: string(fileContents)})
	}

	return jotFiles, nil
}

func captureEvents(app *tview.Application, textArea *tview.TextArea) {

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyESC:
            textArea.SetDisabled(true)
        case tcell.KeyRune:
            switch event.Rune(){
                case 'q':
                    app.Stop()
                case 'i':
                    textArea.SetDisabled(false)
            case 'j':
                // textArea
            case 'k':
                // textArea.

        }
    }
		return event
	})
}

func showUI() error {
	app := tview.NewApplication()


	homeDir := os.Getenv("HOME")

	if homeDir == "" {
		return fmt.Errorf("HOME Env Var not set")
	}

	folderPath := filepath.Join(homeDir, ".config/jot/")

	jotFile, _ := readFiles(folderPath)

	textArea := tview.NewTextArea().
		SetPlaceholder("Enter text here...")
	textArea.SetTitle(jotFile[0].fileName).SetBorder(true)
	textArea.SetText(jotFile[0].fileContents, false)
	textArea.SetDisabled(true)

    captureEvents(app, textArea)

	if err := app.SetRoot(textArea, true).SetFocus(textArea).Run(); err != nil {
		panic(err)
	}
	return nil
}
