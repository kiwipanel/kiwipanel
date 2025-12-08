package provision

import (
	"fmt"
	"path/filepath"
)

// public API

// ProvisionOptions holds all inputs for creating a site.
type ProvisionOptions struct {
	Domain     string // example.com
	SystemUser string // example_com_usr
	PHPVersion string // 8.2
	BasePath   string // /var/www
	Templates  string // path to templates dir
}

// CreateWebsite orchestrates the site creation and returns an error or nil.
// On any failure it attempts a best-effort rollback.
func CreateWebsite(opts ProvisionOptions) error {
	home := filepath.Join(opts.BasePath, opts.SystemUser)
	siteWWW := filepath.Join(home, "data", "www", opts.Domain)

	created := &creationTracker{}

	// 1) Validate
	if !ValidateUsername(opts.SystemUser) {
		return fmt.Errorf("invalid system username: %s", opts.SystemUser)
	}

	// 2) Create system user
	if err := CreateSystemUser(opts.SystemUser, home); err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	created.User = true

	// 3) Create directories and set ownership
	if err := EnsureDirectory(home, opts.SystemUser, 0750); err != nil {
		Rollback(created, opts)
		return fmt.Errorf("create home dir: %w", err)
	}
	created.HomeDir = true

	if err := EnsureDirectory(siteWWW, opts.SystemUser, 0750); err != nil {
		Rollback(created, opts)
		return fmt.Errorf("create webroot: %w", err)
	}
	created.SiteWWW = true

	// 4) write minimal index.php if missing
	if err := WriteDefaultIndex(siteWWW, opts.SystemUser); err != nil {
		Rollback(created, opts)
		return fmt.Errorf("write index: %w", err)
	}

	// 5) Create PHP-FPM pool
	if err := CreatePhpFpmPool(opts.Templates, opts.PHPVersion, opts.SystemUser, home); err != nil {
		Rollback(created, opts)
		return fmt.Errorf("php pool: %w", err)
	}
	created.PhpPool = true

	// 6) Create nginx config
	if err := CreateNginxSite(opts.Templates, opts.Domain, opts.SystemUser, siteWWW); err != nil {
		Rollback(created, opts)
		return fmt.Errorf("nginx site: %w", err)
	}
	created.Nginx = true

	// 7) Reload services: php + nginx
	if err := ReloadPhpFpm(opts.PHPVersion); err != nil {
		Rollback(created, opts)
		return fmt.Errorf("reload php: %w", err)
	}
	if err := ReloadNginx(); err != nil {
		Rollback(created, opts)
		return fmt.Errorf("reload nginx: %w", err)
	}

	return nil
}

// DeleteWebsite deletes configs and optionally system user.
func DeleteWebsite(domain, systemUser, basePath, templatesDir string, removeSystemUser bool) error {
	home := filepath.Join(basePath, systemUser)
	siteWWW := filepath.Join(home, "data", "www", domain)

	// remove nginx site
	if err := RemoveNginxSite(domain); err != nil {
		return err
	}
	// remove php pool
	if err := RemovePhpFpmPool(templatesDir, systemUser); err != nil {
		return err
	}
	// reload services
	_ = ReloadPhpFpmAuto() // try best-effort
	_ = ReloadNginx()

	// optionally remove system user (and home)
	if removeSystemUser {
		if err := DeleteSystemUser(systemUser); err != nil {
			return err
		}
	}

	// best-effort: remove webroot dir if exists
	_ = RemovePath(siteWWW)
	return nil
}
