package provision

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateNginxSite renders nginx template and enables it via symlink.
func CreateNginxSite(templatesDir, domain, systemUser, root string) error {
	avail := "/etc/nginx/sites-available"
	enabled := "/etc/nginx/sites-enabled"
	dest := filepath.Join(avail, domain+".conf")
	data := map[string]string{
		"Domain": domain,
		"User":   systemUser,
		"Root":   root,
	}
	if err := RenderTemplate(templatesDir, "site.tmpl", data, dest); err != nil {
		return err
	}
	// create symlink
	link := filepath.Join(enabled, domain+".conf")
	_ = os.Remove(link)
	if err := os.Symlink(dest, link); err != nil {
		return fmt.Errorf("symlink enable failed: %w", err)
	}
	// test nginx config
	if out, err := ExecCombinedOutput("nginx", "-t"); err != nil {
		return fmt.Errorf("nginx -t failed: %s (%w)", out, err)
	}
	return nil
}

func RemoveNginxSite(domain string) error {
	avail := "/etc/nginx/sites-available"
	enabled := "/etc/nginx/sites-enabled"
	dest := filepath.Join(avail, domain+".conf")
	link := filepath.Join(enabled, domain+".conf")
	_ = os.Remove(link)
	_ = os.Remove(dest)
	return nil
}

func ReloadNginx() error {
	if _, err := ExecCombinedOutput("systemctl", "is-enabled", "--quiet", "nginx"); err != nil {
		// nginx not enabled; attempt safe reload anyway
		if out, rer := ExecCombinedOutput("nginx", "-t"); rer != nil {
			return fmt.Errorf("nginx not enabled and test failed: %s", out)
		}
		return nil
	}
	if out, err := ExecCombinedOutput("nginx", "-t"); err != nil {
		return fmt.Errorf("nginx test failed: %s", out)
	}
	if out, err := ExecCombinedOutput("systemctl", "reload", "nginx"); err != nil {
		return fmt.Errorf("reload nginx failed: %s", out)
	}
	return nil
}
