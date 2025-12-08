package provision

import (
	"path/filepath"
)

type creationTracker struct {
	User    bool
	HomeDir bool
	SiteWWW bool
	PhpPool bool
	Nginx   bool
}

func Rollback(created *creationTracker, opts ProvisionOptions) {
	// best-effort rollback
	home := filepath.Join(opts.BasePath, opts.SystemUser)
	siteWWW := filepath.Join(home, "data", "www", opts.Domain)

	if created.Nginx {
		_ = RemoveNginxSite(opts.Domain)
	}
	if created.PhpPool {
		_ = RemovePhpFpmPool(opts.Templates, opts.SystemUser)
	}
	if created.SiteWWW {
		_ = RemovePath(siteWWW)
	}
	if created.HomeDir {
		// if HomeDir created but user exists: attempt to remove dir
		_ = RemovePath(home)
	}
	if created.User {
		_ = DeleteSystemUser(opts.SystemUser)
	}
}
