package helpers

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	isInstalledOnce sync.Once
	isInstalledVal  bool
)

func IsInstalled() bool {
	isInstalledOnce.Do(func() {
		// method 1: check config exists
		if _, err := os.Stat("/opt/kiwipanel/config/kiwipanel.toml"); err == nil {
			isInstalledVal = true
			return
		}
		// method 2: fallback to binary location
		exe, _ := os.Executable()
		real, _ := filepath.EvalSymlinks(exe)
		isInstalledVal = strings.HasPrefix(real, "/opt/kiwipanel")
	})
	return isInstalledVal
}

func MaintenanceFlagPath() string {
	if IsInstalled() {
		return "/opt/kiwipanel/meta/hidden"
	}
	return "kiwipanel/meta/hidden"
}
