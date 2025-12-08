package provision

import (
	"fmt"
	"os/user"
	"strconv"
)

// ValidateUsername: only allow a-z0-9_ and underscore. Adjust per policy.
func ValidateUsername(s string) bool {
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '_') {
			return false
		}
	}
	return len(s) > 1 && len(s) <= 32
}

// CreateSystemUser creates user with nologin shell and home.
// If user exists, returns nil.
func CreateSystemUser(username, home string) error {
	if UserExists(username) {
		return nil
	}
	// useradd --create-home --home-dir <home> --shell /usr/sbin/nologin <user>
	cmd := ExecCommand("useradd", "--create-home", "--home-dir", home, "--shell", "/usr/sbin/nologin", username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("useradd failed: %w", err)
	}
	return nil
}

// DeleteSystemUser deletes user and its home (-r)
func DeleteSystemUser(username string) error {
	if !UserExists(username) {
		return nil
	}
	cmd := ExecCommand("userdel", "-r", username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("userdel failed: %w", err)
	}
	return nil
}

// UserExists checks /etc/passwd via os/user
func UserExists(username string) bool {
	_, err := user.Lookup(username)
	return err == nil
}

// Helper: get uid/gid; returns -1 on error
func LookupUID(username string) int {
	u, err := user.Lookup(username)
	if err != nil {
		return -1
	}
	uid, _ := strconv.Atoi(u.Uid)
	return uid
}
func LookupGID(username string) int {
	u, err := user.Lookup(username)
	if err != nil {
		return -1
	}
	gid, _ := strconv.Atoi(u.Gid)
	return gid
}
