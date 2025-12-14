package helpers

import "os/exec"

func IsProcessRunning(name string) bool {
	cmd := exec.Command("pgrep", name)
	if err := cmd.Run(); err != nil {
		return false // not running
	}
	return true
}
