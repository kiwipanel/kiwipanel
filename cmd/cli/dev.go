package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kiwipanel/kiwipanel/internal/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dev)
}

var dev = &cobra.Command{
	Use:   "dev",
	Short: "Run development server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running development server...")
		if isInstalled() {
			fmt.Println("❌ Dev mode is disabled on installed servers.")
			fmt.Println("👉 Use `kiwipanel start` instead.")
			os.Exit(1)
		}
		app.Boot("development")
		return nil
	},
}

func isInstalled() bool {
	// method 1: check config exists
	if _, err := os.Stat("/opt/kiwipanel/config/kiwipanel.toml"); err == nil {
		return true
	}
	// method 2: fallback to binary location
	exe, _ := os.Executable()
	real, _ := filepath.EvalSymlinks(exe)
	return strings.HasPrefix(real, "/opt/kiwipanel")
}
