package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update KiwiPanel safely (A/B slots)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("this command must be run as root")
		}
		script := "/opt/kiwipanel/scripts/update.sh"
		if _, err := os.Stat(script); err != nil {
			return fmt.Errorf("update script not found: %s", script)
		}

		color.Green("Starting KiwiPanel update...")
		c := exec.Command(script, "update")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		return c.Run()
	},
}
