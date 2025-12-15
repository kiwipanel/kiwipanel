package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(status)
}

var status = &cobra.Command{
	Use:   "status",
	Short: "Check the status of KiwiPanel",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Check the status of KiwiPanel:")
		return nil
	},
}
