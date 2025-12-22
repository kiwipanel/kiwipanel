package cli

import (
	"fmt"
	"os"

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
}

func Root() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
