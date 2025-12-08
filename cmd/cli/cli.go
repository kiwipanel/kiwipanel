package cli

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(cliCmd)
}

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Utility tools",
	Run: func(cmd *cobra.Command, args []string) {
		println("This is CLI tool")
	},
}
