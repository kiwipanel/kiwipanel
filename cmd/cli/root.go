package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	colorOK    = color.New(color.FgGreen).SprintFunc()
	colorFail  = color.New(color.FgRed).SprintFunc()
	colorWarn  = color.New(color.FgYellow).SprintFunc()
	colorTitle = color.New(color.Bold, color.FgCyan).SprintFunc()
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
