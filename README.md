## KiwiPanel

⚠️ WARNING: PRE-ALPHA RELEASE - DO NOT DEPLOY TO PRODUCTION

KiwiPanel is a lightweight, open-source server control panel focused on simplicity, transparency, and sane defaults. It is designed to help you manage a LOMP stack (Linux, OpenLiteSpeed, MariaDB, PHP) without the bloat, lock-in, or opaque automation commonly found in traditional hosting panels.

Unlike all-in-one panels that attempt to abstract everything away, KiwiPanel aims to stay close to the underlying system. Most operations map directly to standard Linux tools and configurations, making the panel predictable, auditable, and friendly to developers and system administrators who want control rather than magic. 

KiwiPanel is written primarily in Go, with a strong emphasis on:
- Minimal resource usage
- Clear system visibility
- Scriptable and inspectable behavior
- Clean separation between the panel and the server stack

This project is still in early development and evolving rapidly. APIs, features, and internal design may change at any time. Feedback, issue reports, and contributions are welcome.

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

