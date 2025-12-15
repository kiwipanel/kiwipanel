package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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
	rootCmd.AddCommand(checkCmd)
	// checkCmd.Flags().BoolVar(&checkFix, "fix", false, "Attempt safe fixes for common problems")
	// checkCmd.Flags().BoolVar(&checkJSON, "json", false, "Output result as JSON")
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run system checks (services, config, ports, directories)",
	RunE: func(cmd *cobra.Command, args []string) error {
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
	serviceNames := []string{"kiwipanel.service", "openlitespeed", "mariadb", "redis"}
	for _, s := range serviceNames {
		if serviceExists(s) {
			if svcRunning(s) {
				results = append(results, CheckResult{fmt.Sprintf("service:%s", s), true, "running"})
			} else {
				results = append(results, CheckResult{fmt.Sprintf("service:%s", s), false, "stopped"})
				if doFix {
					_ = exec.Command("systemctl", "start", s).Run()
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
		ok, detail := portFree(p)
		if ok {
			results = append(results, CheckResult{fmt.Sprintf("port:%d", p), true, "free"})
		} else {
			results = append(results, CheckResult{fmt.Sprintf("port:%d", p), false, detail})
		}
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

	return &CheckReport{
		Timestamp: time.Now(),
		Results:   results,
	}
}

func printHuman(r *CheckReport) {
	fmt.Println()
	fmt.Println("KiwiPanel System Check")
	fmt.Println("======================")
	for _, it := range r.Results {
		icon := "✖"
		if it.OK {
			icon = "✔"
		}
		fmt.Printf(" [%s] %s - %s\n", icon, padRight(it.Name, 30), it.Details)
	}
	fmt.Printf("\nHealth: %d%% - %s\n", r.HealthPct, r.Summary)
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
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
	// use systemctl if available
	if _, err := exec.LookPath("systemctl"); err == nil {
		err := exec.Command("systemctl", "status", name).Run()
		return err == nil
	}
	// fallback: check if process name exists
	out, _ := exec.Command("ps", "aux").Output()
	return strings.Contains(string(out), name)
}

func svcRunning(name string) bool {
	if _, err := exec.LookPath("systemctl"); err == nil {
		out, _ := exec.Command("systemctl", "is-active", name).Output()
		return strings.TrimSpace(string(out)) == "active"
	}
	return false
}

func portFree(port int) (bool, string) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false, err.Error()
	}
	ln.Close()
	return true, ""
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
