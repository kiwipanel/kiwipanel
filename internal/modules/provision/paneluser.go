package provision

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// Defaults (customize if you like)
const (
	DefaultUser     = "kiwipanel"
	DefaultHome     = "/opt/kiwipanel"
	DefaultShell    = "/usr/sbin/nologin"
	SystemdUnitPath = "/etc/systemd/system/kiwipanel.service"
	SudoersPath     = "/etc/sudoers.d/kiwipanel"
	DefaultLogDir   = "/var/log/kiwipanel"
	AllowedSudoCmds = "/usr/sbin/useradd, /usr/sbin/userdel, /bin/chown, /bin/mkdir, /bin/rm, /bin/ln, /usr/bin/systemctl, /usr/sbin/usermod, /bin/chmod"
)

// EnsurePanelUser ensures the system user and base directories exist.
// It will create the system user (system account), create home, set ownership,
// create log dir, write minimal systemd unit and sudoers snippet (if writeConfig true).
func EnsurePanelUser(username, home string, writeConfig bool) error {
	if username == "" {
		username = DefaultUser
	}
	if home == "" {
		home = DefaultHome
	}

	// Create system user if missing
	if err := createSystemUserIfMissing(username, home); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	// Ensure dirs
	if err := ensureDirOwned(home, username, 0755); err != nil {
		return fmt.Errorf("ensure home dir: %w", err)
	}
	if err := ensureDirOwned(DefaultLogDir, username, 0755); err != nil {
		return fmt.Errorf("ensure log dir: %w", err)
	}

	// create a minimal skeleton inside home
	if err := ensurePanelSkeleton(home, username); err != nil {
		return fmt.Errorf("skeleton: %w", err)
	}

	if writeConfig {
		// write systemd unit (optional - only writes file, does not enable)
		if err := writeSystemdUnit(username, home); err != nil {
			return fmt.Errorf("write systemd unit: %w", err)
		}
		// write sudoers snippet with limited permissions
		if err := writeSudoersSnippet(username); err != nil {
			return fmt.Errorf("write sudoers: %w", err)
		}
	}

	return nil
}

// createSystemUserIfMissing uses useradd to create a system user. If exists, no-op.
func createSystemUserIfMissing(username, home string) error {
	if userExists(username) {
		return nil
	}

	// useradd --system --create-home --home-dir <home> --shell <shell> <username>
	cmd := exec.Command("useradd",
		"--system",
		"--create-home",
		"--home-dir", home,
		"--shell", DefaultShell,
		"--comment", "KiwiPanel system user",
		username,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("useradd failed: %s: %w", strings.TrimSpace(string(out)), err)
	}
	return nil
}

func userExists(username string) bool {
	_, err := user.Lookup(username)
	return err == nil
}

// ensureDirOwned makes directory and chowns to user:group
func ensureDirOwned(path, username string, mode os.FileMode) error {
	if err := os.MkdirAll(path, mode); err != nil {
		return fmt.Errorf("mkdir %s: %w", path, err)
	}
	uid := lookupUID(username)
	gid := lookupGID(username)
	if uid < 0 || gid < 0 {
		return fmt.Errorf("lookup uid/gid failed for %s", username)
	}
	if err := os.Chown(path, uid, gid); err != nil {
		return fmt.Errorf("chown %s: %w", path, err)
	}
	return nil
}

func ensurePanelSkeleton(home, username string) error {
	paths := []string{
		filepath.Join(home, "bin"),
		filepath.Join(home, "templates"),
		filepath.Join(home, "data"),
		filepath.Join(home, "tmp"),
	}
	for _, p := range paths {
		if err := ensureDirOwned(p, username, 0755); err != nil {
			return err
		}
	}
	// ensure a README or marker file
	marker := filepath.Join(home, "README.txt")
	if _, err := os.Stat(marker); os.IsNotExist(err) {
		_ = os.WriteFile(marker, []byte("KiwiPanel system files\n"), 0644)
		uid := lookupUID(username)
		gid := lookupGID(username)
		_ = os.Chown(marker, uid, gid)
	}
	return nil
}

func lookupUID(username string) int {
	u, err := user.Lookup(username)
	if err != nil {
		return -1
	}
	uid, _ := strconv.Atoi(u.Uid)
	return uid
}
func lookupGID(username string) int {
	u, err := user.Lookup(username)
	if err != nil {
		return -1
	}
	gid, _ := strconv.Atoi(u.Gid)
	return gid
}

// writeSystemdUnit writes a recommended systemd unit file for the service.
// It does not enable or start the unit; the operator may run `systemctl daemon-reload && systemctl enable --now kiwipanel`.
func writeSystemdUnit(username, home string) error {
	unit := `[Unit]
Description=KiwiPanel backend service
After=network.target

[Service]
Type=simple
User=` + username + `
Group=` + username + `
WorkingDirectory=` + home + `
ExecStart=` + filepath.Join(home, "bin", "kiwipanel") + ` server
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
`
	if err := os.WriteFile(SystemdUnitPath, []byte(unit), 0644); err != nil {
		return fmt.Errorf("write systemd unit: %w", err)
	}
	return nil
}

// writeSudoersSnippet writes a narrow sudoers file allowing the panel user to run only specific commands with NOPASSWD.
// Note: this file must be owned by root and mode 0440.
func writeSudoersSnippet(panelUser string) error {
	contents := panelUser + " ALL=(ALL) NOPASSWD: " + AllowedSudoCmds + "\n"
	if err := os.WriteFile(SudoersPath, []byte(contents), 0440); err != nil {
		return fmt.Errorf("write sudoers: %w", err)
	}
	// ensure correct ownership (root:root) and perms
	if err := os.Chown(SudoersPath, 0, 0); err != nil {
		return fmt.Errorf("chown sudoers: %w", err)
	}
	if err := os.Chmod(SudoersPath, 0440); err != nil {
		return fmt.Errorf("chmod sudoers: %w", err)
	}
	return nil
}

// RemovePanelUser removes sudoers, systemd unit, and (optionally) the system user and home.
// Careful: run this as root.
func RemovePanelUser(username string, removeUser bool) error {
	// remove sudoers file
	_ = os.Remove(SudoersPath)
	// remove systemd unit
	_ = os.Remove(SystemdUnitPath)
	if removeUser {
		// userdel --remove <username>
		cmd := exec.Command("userdel", "--remove", username)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("userdel failed: %s: %w", string(out), err)
		}
	}
	return nil
}
