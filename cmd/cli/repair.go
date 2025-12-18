package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(repairCmd)
}

var repairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Repair KiwiPanel file ownership, permissions, and service",
	RunE: func(cmd *cobra.Command, args []string) error {
		script := "/opt/kiwipanel/scripts/repair.sh"

		if _, err := os.Stat(script); err != nil {
			return fmt.Errorf("repair script not found: %s", script)
		}

		c := exec.Command(script)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}
