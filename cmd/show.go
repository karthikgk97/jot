package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show notes",
	Long:  `Shows allows you to view notes`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("Executing Show Cmd")       
	},
}

var showTopCmd = &cobra.Command{
	Use:   "top",
	Short: "Show notes from a top level view",
	Long:  `Show top allows you to view the top line contents of the notes.`,
}


var showBottomCmd = &cobra.Command{
	Use:   "bottom",
	Short: "Show notes from a bottom level view",
	Long:  `Show top allows you to view the bottom contents of the notes.`,
}


func init() {
	showCmd.PersistentFlags().StringP("file", "f", "", `
  The file to show.
  Defaults to the one provided in config
    `)
	showCmd.AddCommand(showTopCmd)
	showCmd.AddCommand(showBottomCmd)
}
