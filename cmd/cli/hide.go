package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kiwipanel/kiwipanel/pkg/helpers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(hide)
}

var hide = &cobra.Command{
	Use:   "hide",
	Short: "Turn off the KiwiPanel web interface temporarily",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Turning off KiwiPanel web interface...")
		metaPath := helpers.MaintenanceFlagPath()
		if err := os.WriteFile(metaPath, []byte("maintenance"), 0600); err != nil {
			return fmt.Errorf("Failed to hide KiwiPanel: %w", err)
		}
		message := "KiwiPanel web interface is now hidden (maintenance mode). This mode can be used as a security measure. In case you want to unhide it, run `kiwipanel show` command."
		color.Cyan(message)

		return nil
	},
}
