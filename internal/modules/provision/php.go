package provision

import (
	"fmt"
	"os"
	"path/filepath"
)

func phpPoolPath(phpVersion, systemUser string) string {
	return filepath.Join("/etc/php", phpVersion, "fpm", "pool.d", systemUser+".conf")
}

func CreatePhpFpmPool(templatesDir, phpVersion, systemUser, homeDir string) error {
	dest := phpPoolPath(phpVersion, systemUser)
	data := map[string]string{
		"User": systemUser,
		"Home": homeDir,
		"Pool": systemUser,
	}
	if err := RenderTemplate(templatesDir, "php_pool.tmpl", data, dest); err != nil {
		return err
	}
	return nil
}

func RemovePhpFpmPool(templatesDir, systemUser string) error {
	// try multiple PHP versions? keeping simple: remove common path
	pattern := filepath.Join("/etc/php", "*", "fpm", "pool.d", systemUser+".conf")
	_ = RemoveGlob(pattern) // best-effort
	return nil
}

func ReloadPhpFpm(version string) error {
	service := fmt.Sprintf("php%s-fpm", version)
	if _, err := ExecCombinedOutput("systemctl", "is-enabled", "--quiet", service); err != nil {
		// not enabled; try restart fallback
		if out, err2 := ExecCombinedOutput("systemctl", "restart", service); err2 != nil {
			return fmt.Errorf("php service not enabled and restart failed: %s", out)
		}
		return nil
	}
	if out, err := ExecCombinedOutput("systemctl", "reload", service); err != nil {
		return fmt.Errorf("php reload failed: %s", out)
	}
	return nil
}

// ReloadPhpFpmAuto tries to reload any php-fpm service (best-effort)
func ReloadPhpFpmAuto() error {
	_, _ = ExecCombinedOutput("systemctl", "reload", "php7.4-fpm")
	_, _ = ExecCombinedOutput("systemctl", "reload", "php8.0-fpm")
	_, _ = ExecCombinedOutput("systemctl", "reload", "php8.1-fpm")
	_, _ = ExecCombinedOutput("systemctl", "reload", "php8.2-fpm")
	return nil
}

// RemoveGlob helper
func RemoveGlob(pattern string) error {
	matches, _ := filepath.Glob(pattern)
	for _, f := range matches {
		_ = os.Remove(f)
	}
	return nil
}
