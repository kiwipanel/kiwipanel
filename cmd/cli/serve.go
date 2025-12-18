package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/kiwipanel/kiwipanel/internal/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

// This cleanly separates the systemd service execution from manual CLI usage. We dont show this command to users. It triggers the server when systemd starts the service.
// That is why in the RunE function we check for the INVOCATION_ID environment variable, which systemd sets when starting a service. If it's not set, we return an error telling the user not to run this command manually.
// Additionally, this approach helps to prevent accidental manual execution of the server command, ensuring that the server is only started in the intended environment managed by systemd.
// This is a good practice for services that are meant to be managed by systemd, as it helps maintain the integrity and expected behavior of the service.
// This command is hidden from the help output to avoid confusion for users who might think they should run it manually.
// Here we assume that KiwiPanel is only run as a systemd service in production environments.
// Noted the `serve“ command in /etc/systemd/system/kiwipanel.service  ExecStart=/opt/kiwipanel/bin/kiwipanel serve
var serveCmd = &cobra.Command{
	Use:    "serve",
	Short:  "Run KiwiPanel server (systemd only)",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Optional protection: ensure systemd started us
		if os.Getenv("INVOCATION_ID") == "" {
			return fmt.Errorf("do not run serve manually; use systemctl")
		}
		app.Boot("production")
		log.Println("KiwiPanel server starting (systemd)")
		select {} // KEEP PROCESS ALIVE
		// return nil
	},
}
