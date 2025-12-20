## KiwiPanel

⚠️ WARNING: PRE-ALPHA RELEASE - DO NOT DEPLOY TO PRODUCTION

KiwiPanel is a lightweight, open-source server control panel focused on simplicity, transparency, and sane defaults. It is designed to help you manage a LOMP stack (Linux, OpenLiteSpeed, MariaDB, PHP) without the bloat, lock-in, or opaque automation commonly found in traditional hosting panels.

Unlike all-in-one panels that attempt to abstract everything away, KiwiPanel aims to stay close to the underlying system. Most operations map directly to standard Linux tools and configurations, making the panel predictable, auditable, and friendly to developers and system administrators who want control rather than magic.

KiwiPanel is written primarily in Go, with a strong emphasis on:
- Minimal resource usage
- Clear system visibility
- Scriptable and inspectable behavior
- Clean separation between the panel and the server stack

### Built for Developers & VPS Users

KiwiPanel is designed first and foremost for:
- Developers who want full visibility into how their server works  
- VPS users who prefer lightweight tooling over heavy abstractions  
- System administrators who value reproducibility, auditability, and control  

The panel does **not** attempt to hide Linux behind layers of automation. Instead, it provides a clear interface on top of standard system components, allowing you to learn, debug, and customize your server with confidence.

### No Lock-In, Ever

A core design principle of KiwiPanel is **zero lock-in**:

- KiwiPanel can be **safely uninstalled at any time**
- Your web stack continues to function normally after removal
- No proprietary configuration formats
- No background agents required for runtime operation
- All changes are applied using standard Linux configs and services

KiwiPanel does not replace your system — it works *with* it.

### Practical Server Tasks, Done Transparently

KiwiPanel aims to assist with common server administration tasks while keeping everything explicit and inspectable, including:

- Website and virtual host management
- Database and user management
- Backup primitives (files and databases)
- Security hardening assistance (opt-in, non-destructive)
- Service isolation and sane filesystem layouts
- Log inspection and service health visibility
- TLS and certificate management
- Environment and system diagnostics

Each feature is designed to be:
- Understandable
- Reversible
- Traceable to the underlying system behavior

### Early Development Notice

This project is still in **early development** and evolving rapidly. APIs, features, and internal design may change at any time.

- Expect breaking changes
- Expect missing features
- Expect rough edges

Feedback, issue reports, and contributions are welcome and encouraged. KiwiPanel is built in the open, with the goal of growing into a dependable, no-nonsense control panel that respects both the server and the user.

### Install

**Option 1:**
```bash
bash <(wget -qO- https://raw.githubusercontent.com/kiwipanel/kiwipanel/main/install)
```

**Option 2:**
```bash
curl -sLO https://raw.githubusercontent.com/kiwipanel/kiwipanel/main/install && chmod +x install && sudo bash install
```

##### Port 8443: 
On some cloud service providers such as Amazon Lightsail or Oracle, you have to manually open the port 8443 inside their control dashboards.

##### Supported Operating Systems:
Kiwipanel supports the following operating systems, given that OpenLiteSpeed supports current and non-EOL versions of the following Linux distributions:
- CentOS* 8, 9, 10
- Debian  11, 12, 13
- Ubuntu  22, 24

*Includes RedHat Enterprise Linux and derivatives: AlmaLinux, CloudLinux, Oracle Linux, RockyLinux, VzLinux, etc.

## Roadmap

KiwiPanel is intentionally developed in small, auditable steps. The roadmap below reflects the current direction, but priorities may shift based on real-world usage and community feedback.

### Phase 0 — Foundation (Current / Pre-Alpha)
**Goal:** Establish a clean, inspectable core with minimal abstraction.

- [x] Installer bootstrap for supported Linux distributions  
- [x] Go-based backend architecture  
- [x] CLI framework (`kiwipanel`) for system-level operations  
- [x] SQLite-based local state (no external dependencies)  
- [x] Basic system inspection (CPU, memory, disk, OS)  
- [x] OpenLiteSpeed + MariaDB + PHP stack provisioning  
- [x] Clear separation between panel logic and system tooling  
- [ ] Internal logging and structured error handling (in progress)

⚠️ Not production-ready. Breaking changes expected.

---

### Phase 1 — Core Panel Features (Alpha)
**Goal:** Make KiwiPanel usable for real servers with limited scope.

- [ ] Web UI authentication (local users only)
- [ ] Service management (start/stop/reload):
  - OpenLiteSpeed
  - MariaDB
  - PHP
- [ ] Log viewer (OLS, MariaDB, system logs)
- [ ] Basic website management:
  - Virtual hosts
  - Document roots
  - PHP version selection
- [ ] Database management (create/delete users & databases)
- [ ] Safe defaults for permissions and filesystem layout
- [ ] Non-destructive config generation (no silent overwrites)

---

### Phase 2 — Security & Hardening (Beta)
**Goal:** Secure-by-default without hiding the system.

- [ ] First-run security audit
- [ ] Firewall awareness (UFW / firewalld detection, not replacement)
- [ ] SSH hardening recommendations (opt-in)
- [ ] Fail2ban integration (optional, transparent configs)
- [ ] TLS management (Let’s Encrypt)
- [ ] Port exposure and service visibility
- [ ] Explicit warnings for unsafe configurations

---

### Phase 3 — Developer & Ops Experience
**Goal:** Make KiwiPanel friendly for developers and sysadmins.

- [ ] Full CLI parity with Web UI
- [ ] Scriptable actions (JSON / stdout-friendly output)
- [ ] Backup primitives (files + databases)
- [ ] Restore workflows
- [ ] Environment inspection & health checks
- [ ] Versioned configuration snapshots
- [ ] Read-only / audit mode

---

### Phase 4 — Extensibility & Ecosystem
**Goal:** Grow without becoming bloated.

- [ ] Plugin / module system (strictly sandboxed)
- [ ] Hook system for install/update events
- [ ] External monitoring integration (Prometheus-compatible metrics)
- [ ] API stabilization
- [ ] Documentation-first extensibility guidelines

---

### Explicit Non-Goals
To keep KiwiPanel focused, the following are **intentionally out of scope**:

- ❌ Reseller / multi-tenant hosting abstraction  
- ❌ Proprietary cloud lock-in  
- ❌ Hidden background automation  
- ❌ One-click “magic” that obscures system state  

---

### Philosophy Going Forward
KiwiPanel prioritizes:
- Predictability over convenience  
- Visibility over abstraction  
- Unix principles over panel-driven orchestration  

If a feature cannot be explained clearly or mapped directly to system behavior, it likely does not belong in KiwiPanel.

### Dependencies
- https://github.com/go-chi/chi
- https://github.com/shirou/gopsutil/
- https://github.com/spf13/cobra
- https://github.com/bitfield/script
- https://gorm.io/driver/sqlite
- https://gorm.io/gorm
- https://github.com/ajaxorg/ace
- https://codemirror.net/5/doc/manual.html
- https://echarts.apache.org/
- https://github.com/google/ngx_brotli

### References
- https://github.com/imthenachoman/How-To-Secure-A-Linux-Server
- https://archive.is/8CRGS
- https://docs.openlitespeed.org/installation/script/
- https://github.com/litespeedtech/ols1clk/blob/master/ols1clk.sh
- https://github.com/vernu/vps-audit
- https://github.com/hestiacp/hestiacp
- https://github.com/ConvoyPanel/panel
- https://github.com/wptangtoc/wptangtoc-ols/ or https://archive.is/LJj1a
- https://github.com/QROkes/webinoly
- https://installer.cloudpanel.io/ce/v2/install.sh
- https://github.com/sanvu88/ubuntu-lemp-stack
- https://github.com/duy13/HocVPS-Script/tree/master
- https://github.com/itvn9online/vpssim-free/
- https://github.com/usmannasir/cyberpanel
- https://docs.cloudron.io/

