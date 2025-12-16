## KiwiPanel scaffolding

ALERT: WARNING: PRE-ALPHA RELEASE. DO NOT DEPLOY TO PRODUCTION.

KiwiPanel is a lightweight control panel designed to help you efficiently manage your LEMP (Linux, Nginx, MariaDB, PHP) and LOMP (Linux, OpenLiteSpeed, MariaDB, PHP) stack with minimal hassle.

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

