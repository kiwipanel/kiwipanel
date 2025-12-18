package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check KiwiPanel system health",
	RunE: func(cmd *cobra.Command, args []string) error {
		script := "/opt/kiwipanel/scripts/doctor.sh"

		if _, err := os.Stat(script); err != nil {
			return fmt.Errorf("doctor script not found: %s", script)
		}

		c := exec.Command(script)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		return c.Run()
	},
}
