package helpers

import "os/exec"

func IsServiceRunning(name string) bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", name)
	return cmd.Run() == nil
}
