package web

import (
	"fmt"

	"github.com/kiwipanel/kiwpanel/pkg/bootstrap"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dev)
}

var dev = &cobra.Command{
	Use:   "dev",
	Short: "dev - Start running the app using this command",
	Long:  `The application will be run at the port 8443. The port is configed at the .env`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dev evironment - Start the server using cmd/web/start.go")
		devstart()
	},
}

func devstart() {
	flag := "development"
	bootstrap.Root(flag)
}
