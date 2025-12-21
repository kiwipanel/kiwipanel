package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	logFollow bool
	logLines  int
	logSince  string
	logUntil  string
)

func init() {
	logCmd.Flags().BoolVarP(&logFollow, "follow", "f", false, "Follow log output")
	logCmd.Flags().IntVarP(&logLines, "lines", "n", 100, "Number of lines to show")
	logCmd.Flags().StringVar(&logSince, "since", "", "Show logs since time (e.g. 1h, today)")
	logCmd.Flags().StringVar(&logUntil, "until", "", "Show logs until time")
	rootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "KiwiPanel logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		color.Cyan("Kiwipanel logs...")
		if os.Geteuid() != 0 {
			color.Yellow("⚠ Limited log access (run as root for full logs)")
		}

		args = buildJournalctlArgs()
		color.Cyan("▶ journalctl %v", args)

		cmdExec := exec.Command("journalctl", args...)
		cmdExec.Stdout = os.Stdout
		cmdExec.Stderr = os.Stderr
		cmdExec.Stdin = os.Stdin

		if err := cmdExec.Run(); err != nil {
			return fmt.Errorf("failed to read logs: %w", err)
		}
		return nil
	},
}

func buildJournalctlArgs() []string {
	args := []string{"-u", "kiwipanel"}

	if logFollow {
		args = append(args, "-f")
	}
	if logLines > 0 {
		args = append(args, "-n", strconv.Itoa(logLines))
	}
	if logSince != "" {
		args = append(args, "--since", logSince)
	}
	if logUntil != "" {
		args = append(args, "--until", logUntil)
	}

	return args
}
