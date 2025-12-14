package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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
		fmt.Println("Uninstall KiwiPanel...")

		if os.Geteuid() != 0 {
			return fmt.Errorf("must be run as root")
		}

		uninstallCmd()
		return nil
	},
}

func uninstallCmd() {
	// Try to stop via systemd first
	run("systemctl", "stop", "kiwipanel")
	run("systemctl", "disable", "kiwipanel")

	// Force kill any remaining kiwipanel processes
	run("pkill", "-9", "kiwipanel")

	// Give it a moment to terminate
	time.Sleep(1 * time.Second)

	// Remove the systemd service file
	os.Remove("/etc/systemd/system/kiwipanel.service")
	run("systemctl", "daemon-reload")

	// Now remove the binary
	os.Remove("/usr/local/bin/kiwipanel")

	if !keepData {
		os.RemoveAll("/opt/kiwipanel")
		os.RemoveAll("/var/log/kiwipanel")
	}

	run("userdel", "kiwipanel")
	run("groupdel", "kiwisecure")

	fmt.Println("✅ KiwiPanel successfully uninstalled.")
}

func run(name string, args ...string) {
	fmt.Printf("→ %s %s\n", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  ignored error: %v\n", err)
	}
}
