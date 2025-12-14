package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	run("systemctl stop kiwipanel")
	run("systemctl disable kiwipanel")
	os.Remove("/etc/systemd/system/kiwipanel.service")
	run("systemctl daemon-reload")

	os.Remove("/usr/local/bin/kiwipanel")

	if !keepData {
		os.RemoveAll("/opt/kiwipanel")
		os.RemoveAll("/var/log/kiwipanel")
	}

	run("userdel kiwipanel")
	run("groupdel kiwisecure")

	fmt.Println("✅ KiwiPanel successfully uninstalled.")
}

func run(command string) error {
	fmt.Printf("→ %s\n", command)

	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %s: %w", command, err)
	}

	return nil
}
