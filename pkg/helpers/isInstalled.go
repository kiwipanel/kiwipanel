package helpers

import (
	"os"
	"path/filepath"
	"strings"
)

func IsInstalled() bool {
	// method 1: check config exists
	if _, err := os.Stat("/opt/kiwipanel/config/kiwipanel.toml"); err == nil {
		return true
	}
	// method 2: fallback to binary location
	exe, _ := os.Executable()
	real, _ := filepath.EvalSymlinks(exe)
	return strings.HasPrefix(real, "/opt/kiwipanel")
}
