package cli

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rotatePasscodeCmd)
}

var rotatePasscodeCmd = &cobra.Command{
	Use:   "passcode-rotate",
	Short: "Rotate KiwiPanel passcode",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Rotating KiwiPanel passcode...")

		if os.Geteuid() != 0 {
			fmt.Println("❌ Must be run as root")
			os.Exit(1)
		}

		return rotatePasscode()
	},
}

func confirmRotate() bool {
	fmt.Print("⚠ This will invalidate the current login URL. Continue? [y/N]: ")
	var ans string
	fmt.Scanln(&ans)
	ans = strings.ToLower(strings.TrimSpace(ans))
	return ans == "y" || ans == "yes"
}

func rotatePasscode() error {
	const (
		metaDir      = "/opt/kiwipanel/meta"
		passcodeFile = metaDir + "/passcode"
	)

	if err := os.MkdirAll(metaDir, 0750); err != nil {
		return fmt.Errorf("failed to create meta dir: %w", err)
	}

	// Generate secure random passcode (12 hex chars)
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Errorf("failed to generate passcode: %w", err)
	}

	newPasscode := hex.EncodeToString(buf)

	if err := os.WriteFile(passcodeFile, []byte(newPasscode+"\n"), 0640); err != nil {
		return fmt.Errorf("failed to write passcode: %w", err)
	}

	// Best effort ownership fix
	_ = exec.Command("chown", "kiwipanel:kiwisecure", passcodeFile).Run()

	// Restart service (non-fatal if missing)
	_ = exec.Command("systemctl", "restart", "kiwipanel").Run()

	ip := detectIP()

	fmt.Println("✔ Passcode rotated successfully")
	fmt.Println("New KiwiPanel passcode:", newPasscode)
	if ip != "" {
		fmt.Printf("Login URL: http://%s:8443/%s\n", ip, newPasscode)
	}

	return nil
}

func detectIP() string {
	out, err := exec.Command("hostname", "-I").Output()
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}
