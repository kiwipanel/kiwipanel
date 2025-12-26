package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/kiwipanel/kiwipanel/pkg/helpers"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

type CheckResult struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Details string `json:"details,omitempty"`
}

type CheckReport struct {
	Timestamp time.Time     `json:"timestamp"`
	Results   []CheckResult `json:"results"`
	HealthPct int           `json:"health_pct"`
	Summary   string        `json:"summary"`
}

var (
	checkFix  bool
	checkJSON bool
)

func init() {
	if os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run system checks (services, config, ports, directories, security)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println()
		fmt.Println(colorTitle("KiwiPanel System Check"))
		fmt.Println(colorTitle("======================"))
		fmt.Println(colorOK("→ Running system checks, please wait...\n"))

		helpers.GetUserInfo()

		report := runAllChecks(checkFix)

		// compute health percentage
		okCount := 0
		for _, r := range report.Results {
			if r.OK {
				okCount++
			}
		}
		report.HealthPct = int(float64(okCount) / float64(len(report.Results)) * 100)

		// summary
		if report.HealthPct == 100 {
			report.Summary = "All checks passed"
		} else {
			report.Summary = fmt.Sprintf("%d/%d checks passed", okCount, len(report.Results))
		}

		// output
		if checkJSON {
			b, _ := json.MarshalIndent(report, "", "  ")
			fmt.Println(string(b))
			return nil
		}

		printHuman(report)

		// If critical failures, return error for automation
		if report.HealthPct < 50 {
			return errors.New("critical failures detected (health < 50%)")
		}

		return nil
	},
}

func CheckSSHRootLogin() (bool, string) {
	ok, out := helpers.CmdWithTimeout(5*time.Second, "sshd", "-T")
	if !ok {
		return false, "cannot read sshd config"
	}

	if strings.Contains(out, "permitrootlogin no") || strings.Contains(out, "permitrootlogin prohibit-password") {
		return true, "root login disabled"
	}

	return false, "root login enabled"
}

func CheckSSHPasswordAuth() (bool, string) {
	ok, out := helpers.CmdWithTimeout(5*time.Second, "sshd", "-T")
	if !ok {
		return false, "cannot read sshd config"
	}

	if strings.Contains(out, "passwordauthentication no") {
		return true, "password authentication disabled"
	}

	return false, "password authentication enabled"
}

func CheckSSHPort() (bool, string) {
	ok, out := helpers.CmdWithTimeout(5*time.Second, "sshd", "-T")
	if !ok {
		return false, "cannot read sshd config"
	}

	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "port ") {
			port := strings.TrimSpace(strings.TrimPrefix(line, "port "))
			if port == "22" {
				return false, "using default SSH port 22"
			}
			return true, "using custom port " + port
		}
	}

	return false, "SSH port not detected"
}

func CheckFirewall() (bool, string) {
	// Check all firewalls and prioritize active ones

	// 1. Check firewalld
	if _, err := exec.LookPath("firewall-cmd"); err == nil {
		out, _ := exec.Command("firewall-cmd", "--state").Output()
		if strings.TrimSpace(string(out)) == "running" {
			return true, "firewalld running"
		}
	}

	// 2. Check UFW
	if _, err := exec.LookPath("ufw"); err == nil {
		out, _ := exec.Command("ufw", "status").Output()
		if strings.Contains(string(out), "Status: active") {
			return true, "ufw active"
		}
	}

	// 3. Check nftables
	if _, err := exec.LookPath("nft"); err == nil {
		out, _ := exec.Command("nft", "list", "ruleset").Output()
		if strings.Contains(string(out), "table") {
			return true, "nftables rules present"
		}
	}

	// 4. Check iptables (legacy fallback)
	out, err := exec.Command("iptables", "-L", "-n").Output()
	if err == nil && len(out) > 0 {
		// Has some rules beyond default empty chains
		if strings.Count(string(out), "Chain") > 3 {
			return true, "iptables active"
		}
	}

	// If we get here, no active firewall found
	// Report what's installed but inactive
	if _, err := exec.LookPath("ufw"); err == nil {
		return false, "ufw installed but inactive"
	}
	if _, err := exec.LookPath("firewall-cmd"); err == nil {
		return false, "firewalld installed but inactive"
	}

	return false, "no firewall detected"
}

func CheckIntrusionPrevention() (bool, string) {
	if svcRunning("fail2ban") {
		return true, "fail2ban active"
	}
	if svcRunning("crowdsec") {
		return true, "crowdsec active"
	}
	return false, "no intrusion prevention system (fail2ban || crowdsec) running"
}

func CheckDiskUsage() (bool, string) {
	type fsStat struct {
		FS      string
		Mounts  []string
		UsedPct int
		FreeGB  float64
	}

	paths := []string{"/", "/var", "/opt", "/home"}
	fsMap := make(map[string]*fsStat)

	critical := false
	warning := false

	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			continue
		}

		out, err := exec.Command("df", "-P", p).Output()
		if err != nil {
			continue
		}

		lines := strings.Split(string(out), "\n")
		if len(lines) < 2 {
			continue
		}

		fields := strings.Fields(lines[1])
		if len(fields) < 6 {
			continue
		}

		fs := fields[0]
		usedPct := atoiSafe(strings.TrimSuffix(fields[4], "%"))
		availKB := atoiSafe(fields[3])
		freeGB := float64(availKB) / 1024.0 / 1024.0

		stat, ok := fsMap[fs]
		if !ok {
			stat = &fsStat{
				FS:      fs,
				UsedPct: usedPct,
				FreeGB:  freeGB,
			}
			fsMap[fs] = stat
		}

		stat.Mounts = append(stat.Mounts, p)

		if usedPct >= 85 {
			critical = true
		} else if usedPct >= 70 {
			warning = true
		}
	}

	if len(fsMap) == 0 {
		return true, "no relevant filesystems detected"
	}

	var parts []string
	for _, s := range fsMap {
		parts = append(parts, fmt.Sprintf(
			"%s (%s): %d%% used, %.1f GB free",
			s.FS,
			strings.Join(s.Mounts, ","),
			s.UsedPct,
			s.FreeGB,
		))
	}

	summary := strings.Join(parts, " | ")

	if critical {
		return false, "disk critical: " + summary
	}
	if warning {
		return true, "disk warning: " + summary
	}

	return true, "disk OK: " + summary
}

func CheckMemoryUsage() (bool, string) {
	out, err := exec.Command("free").Output()
	if err != nil {
		return true, "cannot check memory"
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "Mem:") {
			f := strings.Fields(line)
			used := atoiSafe(f[2])
			total := atoiSafe(f[1])
			pct := (used * 100) / total

			if pct >= 80 {
				return false, fmt.Sprintf("memory usage %d%%", pct)
			}
			return true, fmt.Sprintf("memory usage %d%%", pct)
		}
	}

	return true, "memory info unavailable"
}

func CheckWorldWritableDirs() (bool, string) {
	paths := []string{"/opt", "/home", "/var/www", "/etc", "/root"}
	out, err := exec.Command("find", append(paths, "-xdev", "-type", "d", "-perm", "-0002")...).Output()
	if err != nil {
		return true, "cannot scan directories"
	}

	var bad []string
	for _, d := range strings.Split(string(out), "\n") {
		if d != "" && !strings.HasPrefix(d, "/tmp") && !strings.HasPrefix(d, "/var/tmp") {
			bad = append(bad, d)
		}
	}

	if len(bad) > 0 {
		return false, fmt.Sprintf("%d world-writable dirs", len(bad))
	}

	return true, "no unsafe world-writable dirs"
}

func CheckSuspiciousSUID() (bool, string) {
	paths := []string{"/opt", "/home", "/usr/local", "/root"}
	out, err := exec.Command("find", append(paths, "-xdev", "-perm", "-4000", "-type", "f")...).Output()
	if err != nil {
		return true, "cannot scan suid files"
	}

	known := []string{"sudo", "passwd", "su", "mount", "umount"}
	var suspicious int

	for _, f := range strings.Split(string(out), "\n") {
		safe := false
		for _, k := range known {
			if strings.Contains(f, k) {
				safe = true
				break
			}
		}
		if f != "" && !safe {
			suspicious++
		}
	}

	if suspicious > 0 {
		return false, fmt.Sprintf("%d suspicious suid files", suspicious)
	}

	return true, "no suspicious suid files"
}

func runAllChecks(doFix bool) *CheckReport {
	results := []CheckResult{}

	// 1. Binary check
	if path, err := exec.LookPath("kiwipanel"); err == nil {
		results = append(results, CheckResult{"binary", true, fmt.Sprintf("found at %s", path)})
	} else {
		results = append(results, CheckResult{"binary", false, "kiwipanel not in PATH"})
		if doFix {
			// Attempt to link if possible
			_ = safeCreateSymlink("/opt/kiwipanel/bin/kiwipanel", "/usr/local/bin/kiwipanel")
			// re-check
			if path2, err2 := exec.LookPath("kiwipanel"); err2 == nil {
				results = append(results, CheckResult{"binary_fix", true, fmt.Sprintf("linked -> %s", path2)})
			}
		}
	}

	// 2. Config check
	cfgPaths := []string{"/opt/kiwipanel/config/kiwipanel.toml"}
	cfgFound := false
	for _, p := range cfgPaths {
		if _, err := os.Stat(p); err == nil {
			cfgFound = true
			if ok, detail := validateToml(p); ok {
				results = append(results, CheckResult{"config", true, p})
			} else {
				results = append(results, CheckResult{"config", false, detail})
				if doFix {
					// no automatic fix for config (dangerous) — but we can backup and create skeleton
					backupPath := p + ".bak"
					_ = os.Rename(p, backupPath)
					_ = createSkeletonConfig(p)
					results = append(results, CheckResult{"config_fix", true, "backed up and created skeleton: " + backupPath})
				}
			}
			break
		}
	}
	if !cfgFound {
		results = append(results, CheckResult{"config", false, "no config found in /opt or /etc"})
	}

	// 3. Services
	serviceNames := []string{"kiwipanel", "lsws", "mariadb", "redis", "nginx", "apache2", "httpd"}
	for _, s := range serviceNames {
		if serviceExists(s) { // serviceExists will add .service suffix internally
			if svcRunning(s) {
				results = append(results, CheckResult{fmt.Sprintf("service:%s", s), true, "running"})
			} else {
				results = append(results, CheckResult{fmt.Sprintf("service:%s", s), false, "stopped"})
				if doFix {
					serviceName := s
					if !strings.HasSuffix(serviceName, ".service") {
						serviceName = serviceName + ".service"
					}
					_ = exec.Command("systemctl", "start", serviceName).Run()
					if svcRunning(s) {
						results = append(results, CheckResult{fmt.Sprintf("service_fix:%s", s), true, "started"})
					}
				}
			}
		} else {
			results = append(results, CheckResult{fmt.Sprintf("service:%s", s), false, "not installed"})
		}
	}

	// 4. Ports
	ports := []int{80, 443, 7080, 8443, 3306, 6379}
	for _, p := range ports {
		ok, detail := checkPort(p)

		results = append(results, CheckResult{
			Name:    fmt.Sprintf("port:%d", p),
			OK:      ok,
			Details: detail,
		})
	}

	// 5. Directories
	dirs := []string{"/opt/kiwipanel/bin", "/opt/kiwipanel/config", "/opt/kiwipanel/data", "/opt/kiwipanel/logs", "/var/log/kiwipanel"}
	for _, d := range dirs {
		if _, err := os.Stat(d); err == nil {
			results = append(results, CheckResult{fmt.Sprintf("dir:%s", d), true, ""})
		} else {
			results = append(results, CheckResult{fmt.Sprintf("dir:%s", d), false, "missing"})
			if doFix {
				_ = os.MkdirAll(d, 0750)
				results = append(results, CheckResult{fmt.Sprintf("dir_fix:%s", d), true, "created"})
			}
		}
	}

	// 6. Security checks
	securityChecks := []struct {
		name string
		fn   func() (bool, string)
	}{
		{"system:public_ip", CheckPublicIP},
		{"system:cpu_model", CheckCPUModel},
		{"system:load_average", CheckLoadAverage},
		{"system:uptime_since", CheckUptimeSince},
		{"system:disk_usage", CheckDiskUsage},
		{"system:memory_usage", CheckMemoryUsage},
		{"system:users", CheckSystemUsers},
		{"security:uid0_users", CheckUID0Users},
		{"security:passwordless_sudo", CheckPasswordlessSudo},
		{"security:interactive_shells", CheckInteractiveShells},
		{"security:unattended_upgrades", CheckUnattendedUpgrades},
		{"security:pending_updates", CheckPendingUpdates},
		{"security:firewall", CheckFirewall},
		{"firewall:ufw", CheckUFW},
		{"firewall:iptables", CheckIPTables},
		{"firewall:nftables", CheckNFTables},
		{"security:reboot_required", CheckRebootRequired},
		{"security:ssh_root_login", CheckSSHRootLogin},
		{"security:ssh_password_auth", CheckSSHPasswordAuth},
		{"security:ssh_port", CheckSSHPort},
		{"security:intrusion_prevention", CheckIntrusionPrevention},
		{"security:world_writable_dirs", CheckWorldWritableDirs},
		{"security:suid_files", CheckSuspiciousSUID},
		{"network:external_connectivity", CheckNetworkConnectivity},
		{"network:dns_resolution", CheckDNSResolution},
		{"kiwipanel:process", CheckKiwiPanelProcess},
		{"kiwipanel:permissions", CheckKiwiPanelPermissions},
		{"kiwipanel:certificate_validity", CheckSSLCertificates},
	}

	for _, sc := range securityChecks {
		ok, detail := sc.fn()
		results = append(results, CheckResult{
			Name:    sc.name,
			OK:      ok,
			Details: detail,
		})
	}

	return &CheckReport{
		Timestamp: time.Now(),
		Results:   results,
	}
}

func printHuman(r *CheckReport) {
	fmt.Println()
	fmt.Println(colorTitle("KiwiPanel System Check"))
	fmt.Println(colorTitle("======================"))

	for _, it := range r.Results {
		icon := colorFail("✖")
		status := colorFail("FAIL")

		if it.OK {
			icon = colorOK("✔")
			status = colorOK("OK")
		}

		name := padRight(it.Name, 30)
		detail := it.Details
		if detail == "" {
			detail = "-"
		}

		fmt.Printf(" [%s] %-4s %s %s\n", icon, status, name, detail)
	}

	fmt.Println()

	healthColor := colorOK
	if r.HealthPct < 80 {
		healthColor = colorWarn
	}
	if r.HealthPct < 50 {
		healthColor = colorFail
	}

	fmt.Printf(
		"Health: %s - %s\n",
		healthColor(fmt.Sprintf("%d%%", r.HealthPct)),
		healthColor(r.Summary),
	)
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func atoiSafe(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func validateToml(path string) (bool, string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return false, "cannot read config: " + err.Error()
	}

	var tmp map[string]interface{}
	if err := toml.Unmarshal(b, &tmp); err != nil {
		return false, "invalid toml: " + err.Error()
	}

	return true, ""
}

func serviceExists(name string) bool {
	if _, err := exec.LookPath("systemctl"); err != nil {
		return false
	}

	// Ensure .service suffix
	if !strings.HasSuffix(name, ".service") {
		name = name + ".service"
	}

	cmd := exec.Command("systemctl", "show", name, "-p", "LoadState")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	// If service exists, LoadState will be "loaded"
	// If not exists, LoadState will be "not-found"
	return !strings.Contains(string(out), "not-found")
}

func svcRunning(name string) bool {
	if _, err := exec.LookPath("systemctl"); err == nil {
		out, _ := exec.Command("systemctl", "is-active", name).Output()
		return strings.TrimSpace(string(out)) == "active"
	}
	return false
}

func checkPort(port int) (bool, string) {
	cmd := exec.Command("ss", "-lntp")
	out, err := cmd.Output()
	if err != nil {
		// If ss is missing, fail open (do not block)
		return true, "cannot inspect (ss not available)"
	}

	lines := strings.Split(string(out), "\n")
	needle := fmt.Sprintf(":%d", port)

	for _, line := range lines {
		if strings.Contains(line, needle) {
			// Example line:
			// LISTEN 0 4096 0.0.0.0:8443 users:(("kiwipanel",pid=1234,fd=7))
			return false, strings.TrimSpace(line)
		}
	}

	return true, "free"
}

func safeCreateSymlink(src, dst string) error {
	// only create if src exists
	if _, err := os.Stat(src); err != nil {
		return err
	}
	_ = os.Remove(dst) // ignore error
	return os.Symlink(src, dst)
}

func createSkeletonConfig(path string) error {
	skel := `
# KiwiPanel skeleton config
[server]
port = 8443
host = "0.0.0.0"
`
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(skel), 0640)
}

func CheckRebootRequired() (bool, string) {
	const rebootFlag = "/var/run/reboot-required"

	_, err := os.Stat(rebootFlag)
	if err == nil {
		return false, "system requires a restart to apply updates"
	}

	if os.IsNotExist(err) {
		return true, "no restart required"
	}

	if errors.Is(err, fs.ErrPermission) {
		return true, "cannot check reboot flag (permission denied)"
	}

	return false, "unable to determine reboot status"
}

func CheckPublicIP() (bool, string) {
	ok, out := helpers.CmdWithTimeout(5*time.Second, "curl", "-s", "--max-time", "5", "https://ifconfig.me/ip")
	if !ok {
		return false, "unable to detect public IP"
	}

	ip := strings.TrimSpace(out)
	if net.ParseIP(ip) == nil {
		return false, "invalid public IP response"
	}

	return true, ip
}

func CheckDNSResolution() (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addrs, err := net.DefaultResolver.LookupHost(ctx, "8.8.8.8")
	if err != nil || len(addrs) == 0 {
		return false, "DNS resolution failed"
	}
	return true, fmt.Sprintf("DNS working (%s)", addrs[0])
}

func CheckCPUModel() (bool, string) {
	b, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return false, "cannot read /proc/cpuinfo"
	}

	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return true, strings.TrimSpace(parts[1])
			}
		}
	}

	return false, "cpu model not found"
}

func CheckLoadAverage() (bool, string) {
	b, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return false, "cannot read load average"
	}

	parts := strings.Fields(string(b))
	if len(parts) < 3 {
		return false, "invalid loadavg format"
	}

	return true, fmt.Sprintf("1m=%s 5m=%s 15m=%s", parts[0], parts[1], parts[2])
}

func CheckUptimeSince() (bool, string) {
	out, err := exec.Command("uptime", "-s").Output()
	if err != nil {
		return false, "cannot determine uptime start"
	}

	return true, strings.TrimSpace(string(out))
}

func CheckUnattendedUpgrades() (bool, string) {
	if _, err := exec.LookPath("dpkg"); err != nil {
		return true, "not a Debian-based system"
	}

	out, _ := exec.Command("dpkg", "-l").Output()
	if strings.Contains(string(out), "unattended-upgrades") {
		return true, "automatic security updates enabled"
	}

	return false, "unattended-upgrades not installed"
}

func CheckPendingUpdates() (bool, string) {
	if _, err := exec.LookPath("apt-get"); err != nil {
		return true, "apt not available"
	}

	out, err := exec.Command("apt-get", "-s", "upgrade").Output()
	if err != nil {
		return true, "cannot simulate upgrades"
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, " upgraded,") {
			fields := strings.Fields(line)
			if len(fields) > 0 && fields[0] != "0" {
				return false, fields[0] + " updates pending"
			}
		}
	}

	return true, "system up to date"
}

func CheckUFW() (bool, string) {
	if _, err := exec.LookPath("ufw"); err != nil {
		return true, "ufw not installed"
	}

	out, _ := exec.Command("ufw", "status").Output()
	if strings.Contains(string(out), "Status: active") {
		return true, "ufw active"
	}

	return false, "ufw installed but inactive"
}

func CheckIPTables() (bool, string) {
	if _, err := exec.LookPath("iptables"); err != nil {
		return true, "iptables not installed"
	}

	out, err := exec.Command("iptables", "-L", "-n").Output()
	if err != nil {
		return false, "cannot read iptables rules"
	}

	if strings.Contains(string(out), "Chain INPUT") {
		return true, "iptables rules present"
	}

	return false, "no iptables rules found"
}

func CheckNFTables() (bool, string) {
	if _, err := exec.LookPath("nft"); err != nil {
		return true, "nftables not installed"
	}

	out, err := exec.Command("nft", "list", "ruleset").Output()
	if err != nil {
		return false, "cannot read nftables ruleset"
	}

	if strings.Contains(string(out), "table") {
		return true, "nftables rules active"
	}

	return false, "nftables installed but empty ruleset"
}

func CheckSystemUsers() (bool, string) {
	b, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return false, "cannot read /etc/passwd"
	}

	total := 0
	login := 0

	for _, line := range strings.Split(string(b), "\n") {
		if line == "" {
			continue
		}

		total++

		parts := strings.Split(line, ":")
		if len(parts) < 7 {
			continue
		}

		uid, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}

		shell := parts[6]

		if uid >= 1000 &&
			!strings.Contains(shell, "nologin") &&
			!strings.Contains(shell, "false") {
			login++
		}
	}

	return true, fmt.Sprintf("total=%d, login-capable=%d", total, login)
}

func CheckUID0Users() (bool, string) {
	b, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return false, "cannot read /etc/passwd"
	}

	var uid0 []string

	for _, line := range strings.Split(string(b), "\n") {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}

		if parts[2] == "0" {
			uid0 = append(uid0, parts[0])
		}
	}

	if len(uid0) == 1 && uid0[0] == "root" {
		return true, "only root has UID 0"
	}

	return false, "UID 0 users: " + strings.Join(uid0, ", ")
}

func CheckPasswordlessSudo() (bool, string) {
	if _, err := exec.LookPath("sudo"); err != nil {
		return true, "sudo not installed"
	}

	cmd := exec.Command("sudo", "-l", "-U", "root")
	out, err := cmd.CombinedOutput()
	if err != nil {
		// cannot list sudo rules → fail open
		return true, "cannot inspect sudo rules"
	}

	if strings.Contains(string(out), "NOPASSWD") {
		return false, "passwordless sudo detected"
	}

	return true, "no passwordless sudo rules"
}

func CheckInteractiveShells() (bool, string) {
	b, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return false, "cannot read /etc/passwd"
	}

	var users []string

	for _, line := range strings.Split(string(b), "\n") {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 7 {
			continue
		}

		user := parts[0]
		uid, _ := strconv.Atoi(parts[2])
		shell := parts[6]

		if uid >= 1000 &&
			!strings.Contains(shell, "nologin") &&
			!strings.Contains(shell, "false") {
			users = append(users, user)
		}
	}

	if len(users) == 0 {
		return true, "no interactive users"
	}

	return true, fmt.Sprintf("%d interactive users", len(users))
}

func CheckNetworkConnectivity() (bool, string) {
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", 5*time.Second)
	if err != nil {
		return false, fmt.Sprintf("no external connectivity: %v", err)
	}
	conn.Close()
	return true, "external connectivity OK (8.8.8.8:53)"
}

func CheckKiwiPanelProcess() (bool, string) {
	out, _ := exec.Command("systemctl", "is-active", "kiwipanel.service").Output()
	if strings.TrimSpace(string(out)) == "active" {
		return true, "kiwipanel service running"
	}
	return false, "kiwipanel service not running"
}

func CheckKiwiPanelPermissions() (bool, string) {
	dirs := map[string]string{
		"/opt/kiwipanel":     "750",
		"/etc/kiwipanel":     "750",
		"/var/log/kiwipanel": "750",
	}

	for path, expectedPerm := range dirs {
		info, err := os.Stat(path)
		if err != nil {
			return false, fmt.Sprintf("%s missing", path)
		}

		actualPerm := fmt.Sprintf("%o", info.Mode().Perm())
		if actualPerm != expectedPerm {
			return false, fmt.Sprintf("%s has %s (expected %s)", path, actualPerm, expectedPerm)
		}
	}
	return true, "all directories have correct permissions"
}

func CheckSSLCertificates() (bool, string) {
	certPath := "/opt/kiwipanel/config/ssl/certificate.pem"
	if _, err := os.Stat(certPath); err != nil {
		return false, "SSL certificate not found"
	}

	// Parse cert and check expiry
	out, _ := exec.Command("openssl", "x509", "-in", certPath, "-noout", "-enddate").Output()
	if len(out) > 0 {
		return true, strings.TrimSpace(strings.TrimPrefix(string(out), "notAfter="))
	}
	return false, "cannot read certificate expiry"
}
