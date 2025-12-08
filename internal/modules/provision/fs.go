package provision

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDirectory creates the directory and chown to the system user.
func EnsureDirectory(path, username string, mode os.FileMode) error {
	if err := os.MkdirAll(path, mode); err != nil {
		return fmt.Errorf("mkdirall %s: %w", path, err)
	}
	uid := LookupUID(username)
	gid := LookupGID(username)
	if uid < 0 || gid < 0 {
		return fmt.Errorf("lookup uid/gid failed for %s", username)
	}
	if err := os.Chown(path, uid, gid); err != nil {
		return fmt.Errorf("chown %s: %w", path, err)
	}
	return nil
}

func WriteDefaultIndex(siteWWW, username string) error {
	path := filepath.Join(siteWWW, "index.php")
	if _, err := os.Stat(path); err == nil {
		return nil // exists
	}
	content := `<?php
echo "It works - served as user: " . get_current_user();
`
	if err := os.WriteFile(path, []byte(content), 0640); err != nil {
		return err
	}
	uid := LookupUID(username)
	gid := LookupGID(username)
	return os.Chown(path, uid, gid)
}

func RemovePath(path string) error {
	return os.RemoveAll(path)
}
