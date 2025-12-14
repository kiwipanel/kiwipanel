package cli

import (
	"fmt"
	"os"

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
		if helpers.IsInstalled() {
			metaPath = hiddenFlagProduction
		} else {
			metaPath = hiddenFlagDev
		}
		if err := os.Remove(metaPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to show KiwiPanel: %w", err)
		}

		fmt.Println("KiwiPanel web interface is now active.")
		return nil
	},
}
