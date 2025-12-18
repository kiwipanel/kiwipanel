package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/kiwipanel/kiwipanel/pkg/helpers"
	"github.com/spf13/cobra"
)

var keepData bool
var force bool

func init() {
	uninstall.Flags().BoolVar(
		&keepData,
		"keep-data",
		false,
		"Keep KiwiPanel data under /opt/kiwipanel",
	)
	uninstall.Flags().BoolVar(
		&force,
		"yes",
		false,
		"Skip confirmation prompt",
	)
	rootCmd.AddCommand(uninstall)
}

var uninstall = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall KiwiPanel from the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("must be run as root")
		}

		if !force && !confirmUninstall() {
			fmt.Println("❌ Uninstall aborted")
			return nil
		}

		fmt.Println("Start uninstalling KiwiPanel...")
		uninstallCmd()
		return nil
	},
}

func confirmUninstall() bool {
	WarningText := colorFail("⚠️  WARNING: KiwiPanel Uninstall:")
	fmt.Println(WarningText)
	fmt.Println("")
	fmt.Println("This action will:")
	fmt.Println("  - Stop and remove KiwiPanel services")
	fmt.Println("  - Remove /etc/kiwipanel")
	fmt.Println("  - Remove system user & group")
	fmt.Println("")
	fmt.Print("Type 'uninstall' to confirm: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = helpers.NormalizeInput(input)

	return input == "uninstall"
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
