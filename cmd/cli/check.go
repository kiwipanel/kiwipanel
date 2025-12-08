package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check command to make sure KiwiPanel works...",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Check command to make sure KiwiPanel works...")
		fmt.Println(".....................")
	},
}
