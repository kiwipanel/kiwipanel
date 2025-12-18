package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start KiwiPanel service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("must be run as root")
		}

		bin := "/opt/kiwipanel/bin/kiwipanel"

		if err := checkBinaryOwnership(bin, "kiwisecure"); err != nil {
			return fmt.Errorf("binary permission check failed: %w", err)
		}

		color.Green("Starting KiwiPanel service...")

		c := exec.Command("systemctl", "start", "kiwipanel.service")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		return c.Run()
	},
}

func checkBinaryOwnership(binPath, expectedGroup string) error {
	info, err := os.Stat(binPath)
	if err != nil {
		return fmt.Errorf("cannot stat %s: %w", binPath, err)
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("unable to read file ownership for %s", binPath)
	}

	// Lookup current owner
	uid := fmt.Sprint(stat.Uid)
	usr, _ := user.LookupId(uid)
	owner := uid
	if usr != nil {
		owner = usr.Username
	}

	// Lookup current group
	gid := fmt.Sprint(stat.Gid)
	grp, _ := user.LookupGroupId(gid)
	group := gid
	if grp != nil {
		group = grp.Name
	}

	// Lookup expected group
	expGrp, err := user.LookupGroup(expectedGroup)
	if err != nil {
		return fmt.Errorf("expected group %q not found on system", expectedGroup)
	}

	// Check executable bit
	if info.Mode()&0111 == 0 {
		return fmt.Errorf(
			"binary is not executable\n  path: %s\n  owner: %s\n  group: %s\n  mode: %o",
			binPath, owner, group, info.Mode().Perm(),
		)
	}

	// Check group mismatch
	if gid != expGrp.Gid {
		return fmt.Errorf(
			"binary group mismatch\n  path: %s\n  owner: %s\n  group: %s\n  expected group: %s",
			binPath, owner, group, expectedGroup,
		)
	}

	// Check parent directory execute permission
	dir := filepath.Dir(binPath)
	dirInfo, err := os.Stat(dir)
	if err == nil && dirInfo.Mode().Perm()&0110 == 0 {
		return fmt.Errorf(
			"no execute permission on binary directory\n  dir: %s\n  mode: %o",
			dir, dirInfo.Mode().Perm(),
		)
	}

	return nil
}
