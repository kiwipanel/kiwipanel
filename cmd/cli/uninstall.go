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
	Short: "Uninstall KiwiPanel from the system",
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
	// 1. Stop and disable systemd service (if exists)
	helpers.Run("systemctl", "stop", "kiwipanel")
	helpers.Run("systemctl", "disable", "kiwipanel")

	// 2. Remove systemd unit
	_ = os.Remove("/etc/systemd/system/kiwipanel.service")
	helpers.Run("systemctl", "daemon-reload")

	// 3. Remove login shell hook
	_ = os.Remove("/etc/profile.d/kiwipanel.sh")

	// 4. Remove binary
	_ = os.Remove("/usr/local/bin/kiwipanel")

	// 5. Remove configuration
	_ = os.RemoveAll("/etc/kiwipanel")

	// 6. Remove runtime data (optional)
	if !keepData {
		_ = os.RemoveAll("/opt/kiwipanel")
		_ = os.RemoveAll("/var/log/kiwipanel")
		fmt.Println("✔ Runtime data removed")
	} else {
		fmt.Println("ℹ Runtime data preserved (--keep-data)")
	}

	// 7. Remove user and group safely
	if helpers.UserExists("kiwipanel") {
		helpers.Run("userdel", "kiwipanel")
	}

	if helpers.GroupExists("kiwisecure") {
		helpers.Run("groupdel", "kiwisecure")
	}

	// 8. Final verification
	if helpers.IsServiceRunning("kiwipanel") {
		fmt.Println("⚠️  KiwiPanel process still running")
	} else {
		fmt.Println("✅ KiwiPanel successfully uninstalled")
	}
}
