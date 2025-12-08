package cli

import (
	"fmt"

	"github.com/kiwipanel/kiwipanel/internal/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dev)
}

var dev = &cobra.Command{
	Use:   "dev",
	Short: "Run development server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running development server...")
		app.Boot("development")
	},
}
