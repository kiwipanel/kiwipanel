package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	colorTitle = color.New(color.FgCyan, color.Bold).SprintFunc()
	colorOK    = color.New(color.FgGreen).SprintFunc()
	colorWarn  = color.New(color.FgYellow).SprintFunc()
	colorFail  = color.New(color.FgRed).SprintFunc()
	colorInfo  = color.New(color.FgCyan).SprintFunc()
	colorBold  = color.New(color.Bold).SprintFunc()
)
var ShortIntroduction = colorTitle("Kiwipanel control CLI")

var rootCmd = &cobra.Command{
	Use:   "kiwipanel",
	Short: ShortIntroduction,
	// PersistentPreRunE runs before ANY subcommand
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Commands that must NEVER be blocked
		switch cmd.Name() {
		case "help", "version", "serve":
			return nil
		}

		return checkRequireRoot()
	},
}

func Root() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// checkRequireRoot ensures the command is run as root (UID 0)
func checkRequireRoot() error {
	if os.Geteuid() != 0 {
		cmdArgs := ""
		if len(os.Args) > 1 {
			cmdArgs = " " + strings.Join(os.Args[1:], " ")
		}
		return fmt.Errorf("this command must be run as root. Please use: sudo kiwipanel%s", cmdArgs)
	}
	return nil
}
