package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kiwipanel/kiwipanel/pkg/helpers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(show)
}

var show = &cobra.Command{
	Use:   "show",
	Short: "Turn on the KiwiPanel web interface",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Turning on KiwiPanel web interface...")
		metaPath := helpers.MaintenanceFlagPath()
		if err := os.Remove(metaPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to show KiwiPanel: %w", err)
		}

		//fmt.Println("KiwiPanel web interface is now active.")

		message := "KiwiPanel web interface is now active. You can access it in your web browser."
		color.Green(message)
		return nil
	},
}
