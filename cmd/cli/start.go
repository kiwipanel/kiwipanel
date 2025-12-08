package cli

import (
	"fmt"

	"github.com/kiwipanel/kiwipanel/internal/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Run production server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running production server...")
		app.Boot("production")
	},
}
