package cli

import (
	"fmt"
	"os"

	"github.com/kiwipanel/kiwipanel/pkg/helpers"
	"github.com/spf13/cobra"
)

var keepData bool

func init() {
	uninstall.Flags().BoolVar(
		&keepData,
		"keep-data",
		false,
		"Keep KiwiPanel data under /opt/kiwipanel",
	)
	rootCmd.AddCommand(uninstall)
}

var uninstall = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall KiwiPanel",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Start uninstalling KiwiPanel...")

		if os.Geteuid() != 0 {
			return fmt.Errorf("Must be helpers.Run as root")
		}
		uninstallCmd()
		return nil
	},
}

func uninstallCmd() {
	// Stop and disable service
	helpers.Run("systemctl", "stop", "kiwipanel")
	helpers.Run("systemctl", "disable", "kiwipanel")

	// Remove systemd unit
	_ = os.Remove("/etc/systemd/system/kiwipanel.service")
	helpers.Run("systemctl", "daemon-reload")

	// Remove binary (safe even while helpers.Running)
	_ = os.RemoveAll("/usr/local/bin/kiwipanel")

	// Remove data
	if !keepData {
		_ = os.RemoveAll("/opt/kiwipanel")
		_ = os.RemoveAll("/var/log/kiwipanel")
	}

	// Remove user/group
	helpers.Run("userdel", "kiwipanel")
	helpers.Run("groupdel", "kiwisecure")

	if helpers.IsProcessRunning("kiwipanel") {
		fmt.Println("⚠️  KiwiPanel process still helpers.Running!")
	} else {
		fmt.Println("✅ No KiwiPanel process found.")
		fmt.Println("✅ KiwiPanel successfully uninstalled.")
	}
}
