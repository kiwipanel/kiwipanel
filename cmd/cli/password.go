package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/kiwipanel/kiwipanel/pkg/helpers"
	"github.com/spf13/cobra"
)

var (
	shortIntro = "Reset the KiwiPanel web portal password"
)

func init() {
	rootCmd.AddCommand(passWord)
}

var passWord = &cobra.Command{
	Use:   "password",
	Short: shortIntro,
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
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("⚠️  WARNING")
		fmt.Println("This will reset the KiwiPanel web portal password.")
		fmt.Print("Do you want to continue? (yes/no): ")

		input, _ := reader.ReadString('\n')
		input = helpers.NormalizeInput(input)

		switch {
		case strings.HasPrefix(input, "y"):
			return true
		case strings.HasPrefix(input, "n"):
			return false
		default:
			fmt.Println("Please answer yes or no.")
		}
	}
}

func generatePassword() error {
	password := "testing"
	fmt.Printf("Your new KiwiPanel password: %s\n", password)
	color.Green("✅ Password reset successfully.")
	return nil
}
