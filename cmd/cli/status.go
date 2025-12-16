package cli

import (
	"bytes"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show KiwiPanel service status",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("systemctl", "is-active", "kiwipanel.service")
		out, err := c.Output()

		status := string(bytes.TrimSpace(out))

		switch status {
		case "active":
			color.Green("● KiwiPanel is running")
		case "inactive":
			color.Yellow("● KiwiPanel is stopped")
		case "failed":
			color.Red("● KiwiPanel has failed")
		default:
			color.White("● KiwiPanel status: %s", status)
		}

		return err
	},
}
