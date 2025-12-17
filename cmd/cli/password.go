package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(passWord)
}

var passWord = &cobra.Command{
	Use:   "password",
	Short: "Reset the KiwiPanel web portal password",
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("❌ must be run as root")
		}

		if !confirmPasswordReset() {
			fmt.Println("Aborted.")
			return nil
		}

		return generatePassword()
	},
}

func confirmPasswordReset() bool {
	fmt.Println("⚠️  WARNING")
	fmt.Println("This will reset the KiwiPanel web portal password.")
	fmt.Print("Do you want to continue? (yes/no): ")

	var input string
	fmt.Scanln(&input)

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "yes" || input == "y"
}

func generatePassword() error {
	fmt.Println("Creating a new password...")
	fmt.Println("KiwiPanel password:", "testing")
	return nil
}
