package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start KiwiPanel service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("must be run as root")
		}

		color.Green("Starting KiwiPanel service...")

		c := exec.Command("systemctl", "start", "kiwipanel.service")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		return c.Run()
	},
}
