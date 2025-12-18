package cli

import (
	"fmt"
	"os"

	"github.com/kiwipanel/kiwipanel/internal/app"
	"github.com/kiwipanel/kiwipanel/pkg/helpers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dev)
}

var dev = &cobra.Command{
	Use:   "dev",
	Short: "Run development server. Does not work on installed servers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running development server...")
		if helpers.IsInstalled() {
			fmt.Println("❌ Dev mode is disabled on installed servers.")
			fmt.Println("👉 Use `kiwipanel start` instead.")
			os.Exit(1)
		}
		app.Boot("development")
		return nil
	},
}
