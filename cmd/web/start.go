package web

import (
	"fmt"

	"github.com/kiwipanel/scaffolding/pkg/bootstrap"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "start",
	Short: "Start running the app using this command",
	Long:  `The application will be run at the port 3003. The port is configed at the .env`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Start the server using cmd/web/start.go")
		start()
	},
}

func start() {
	bootstrap.Root()
}
