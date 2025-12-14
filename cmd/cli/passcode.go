package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(passcode)
}

var passcode = &cobra.Command{
	Use:   "passcode",
	Short: "Showing KiwiPanel passcode...",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Showing KiwiPanel passcode...")
		if os.Geteuid() != 0 {
			fmt.Println("❌ Must be run as root")
			os.Exit(1)
		}
		showPasscode()
		return nil
	},
}

func showPasscode() error {
	data, err := os.ReadFile("/opt/kiwipanel/meta/passcode")
	if err != nil {
		return fmt.Errorf("passcode not found")
	}
	passcode := strings.TrimSpace(string(data))
	fmt.Println("KiwiPanel passcode:", passcode)
	return nil
}
