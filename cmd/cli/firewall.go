package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(firewallCmd)
}

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Show firewall (UFW) rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			fmt.Println("❌ Must be run as root")
			os.Exit(1)
		}

		return showFirewallStatus()
	},
}

func showFirewallStatus() error {
	cmd := exec.Command("ufw", "status")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ufw not available or not enabled")
	}

	lines := strings.Split(out.String(), "\n")

	active := false
	fmt.Println(colorInfo("Firewall status:"), colorOK("active"))
	fmt.Println()

	fmt.Printf("%-10s %-10s %-8s %-20s\n", "PORT", "PROTOCOL", "ACTION", "SOURCE")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Status:") {
			active = strings.Contains(line, "active")
			continue
		}

		// Skip headers / empty
		if line == "" || strings.HasPrefix(line, "To ") || strings.HasPrefix(line, "--") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// Example: 22/tcp ALLOW Anywhere
		portProto := fields[0]
		action := fields[1]
		source := strings.Join(fields[2:], " ")

		pp := strings.Split(portProto, "/")
		port := pp[0]
		proto := ""
		if len(pp) > 1 {
			proto = pp[1]
		}

		fmt.Printf(
			"%-10s %-10s %-8s %-20s\n",
			port,
			proto,
			colorOK(action),
			source,
		)
	}

	if !active {
		fmt.Println(colorWarn("Firewall is not active"))
	}

	return nil
}
