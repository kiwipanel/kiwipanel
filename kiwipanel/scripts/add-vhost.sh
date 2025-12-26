#!/bin/bash
######################################################################
#                  KiwiPanel - Virtual Host Manager                  #
#                                                                    #
#                Author: Vuong Nguyen and contributors               #
#                  Website: https://kiwipanel.org                    #
######################################################################

set -euo pipefail

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
LSWS_DIR="/usr/local/lsws"
LSWS_CTRL="$LSWS_DIR/bin/lswsctrl"
HTTPD_CONF="$LSWS_DIR/conf/httpd_config.conf"
VHOSTS_DIR="$LSWS_DIR/conf/vhosts"
PHP_VER="83"

# Functions
print_green() { echo -e "${GREEN}$1${NC}"; }
print_info() { echo -e "${CYAN}$1${NC}"; }
print_warn() { echo -e "${YELLOW}Warning: $1${NC}"; }
print_error() { echo -e "${RED}Error: $1${NC}"; exit 1; }

# Check root
[[ $EUID -ne 0 ]] && print_error "This script must be run as root"

# Check if OLS is installed
[[ ! -x "$LSWS_CTRL" ]] && print_error "OpenLiteSpeed not found at $LSWS_CTRL"

# Get domain name
if [[ -n "${1:-}" ]]; then
    DOMAIN="$1"
else
    read -p "Enter domain name (e.g., demo.kiwipanel.org): " DOMAIN
fi

# Validate domain
if [[ ! "$DOMAIN" =~ ^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$ ]]; then
    print_error "Invalid domain name: $DOMAIN"
fi

# Sanitize domain for directory/listener names
VHOST_NAME="${DOMAIN//./_}"
VHOST_ROOT="$LSWS_DIR/$DOMAIN"
DOCROOT="$VHOST_ROOT/html"
VHOST_CONF="$VHOSTS_DIR/$VHOST_NAME.conf"
LISTENER_NAME="${VHOST_NAME}_HTTP"

print_info "Creating virtual host for: $DOMAIN"
print_info "Document root: $DOCROOT"
echo ""

# Create directories
mkdir -p "$VHOSTS_DIR"
mkdir -p "$DOCROOT"
mkdir -p "$VHOST_ROOT/logs"

# Create default index.html
cat > "$DOCROOT/index.html" <<EOF
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>$DOMAIN - KiwiPanel</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            max-width: 800px;
            margin: 100px auto;
            padding: 20px;
            text-align: center;
        }
        h1 { color: #2ecc71; }
        .info { 
            background: #f8f9fa;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }
        code {
            background: #e9ecef;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: monospace;
        }
    </style>
</head>
<body>
    <h1>🎉 Virtual Host Active!</h1>
    <p>Domain: <strong>$DOMAIN</strong></p>
    <div class="info">
        <p>This virtual host is powered by <strong>KiwiPanel</strong> + OpenLiteSpeed</p>
        <p>Document root: <code>$DOCROOT</code></p>
    </div>
    <p>Upload your website files to the document root to get started!</p>
</body>
</html>
EOF

# Create PHP info page
cat > "$DOCROOT/info.php" <<'EOF'
<?php
phpinfo();
EOF

# Create virtual host configuration
cat > "$VHOST_CONF" <<EOF
# =========================================================
# Virtual Host: $DOMAIN
# Created: $(date)
# Managed by: KiwiPanel
# =========================================================

docRoot                   $DOCROOT

enableGzip                1

errorlog $VHOST_ROOT/logs/error.log {
  useServer               0
  logLevel                ERROR
  rollingSize             10M
}

accesslog $VHOST_ROOT/logs/access.log {
  useServer               0
  rollingSize             10M
  keepDays                30
}

index  {
  useServer               0
  indexFiles              index.html, index.php
}

scripthandler  {
  add lsapi:lsphp${PHP_VER} php
}

extprocessor lsphp${PHP_VER} {
  type                    lsapi
  address                 uds://tmp/lshttpd/lsphp${PHP_VER}.sock
  maxConns                35
  env                     PHP_LSAPI_CHILDREN=35
  initTimeout             60
  retryTimeout            0
  persistConn             1
  respBuffer              0
  autoStart               1
  path                    $LSWS_DIR/lsphp${PHP_VER}/bin/lsphp
  extMaxIdleTime          300
}

rewrite  {
  enable                  1
  autoLoadHtaccess        1
}

context / {
  location                $DOCROOT
  allowBrowse             1

  rewrite  {
    enable                1
  }
}
EOF

# Backup main config
cp "$HTTPD_CONF" "${HTTPD_CONF}.bak.$(date +%s)"

# Add virtual host to main config if not exists
if ! grep -q "virtualhost $VHOST_NAME {" "$HTTPD_CONF"; then
    cat >> "$HTTPD_CONF" <<EOF

# =========================================================
# Virtual Host: $DOMAIN
# =========================================================
virtualhost $VHOST_NAME {
  vhRoot                  $VHOST_ROOT
  configFile              $VHOST_CONF
  allowSymbolLink         1
  enableScript            1
  restrained              1
}
EOF
    print_green "✓ Virtual host added to httpd_config.conf"
else
    print_info "Virtual host already exists in httpd_config.conf"
fi

# Add listener mapping if not exists
if ! grep -q "map.*$VHOST_NAME.*$DOMAIN" "$HTTPD_CONF"; then
    # Find the HTTP listener block and add mapping
    awk -v vhost="$VHOST_NAME" -v domain="$DOMAIN" '
    /listener HTTP \{/,/^\}/ {
        if (/map/) {
            print
            if (!found) {
                print "    map                     " vhost " " domain
                found=1
            }
            next
        }
    }
    {print}
    ' "$HTTPD_CONF" > "${HTTPD_CONF}.tmp"
    
    mv "${HTTPD_CONF}.tmp" "$HTTPD_CONF"
    print_green "✓ Domain mapping added to HTTP listener"
else
    print_info "Domain mapping already exists"
fi

# Set proper ownership
chown -R lsadm:lsadm "$VHOST_ROOT"
chmod -R 755 "$VHOST_ROOT"
chmod 644 "$VHOST_CONF"

# Graceful restart
print_info "Restarting OpenLiteSpeed..."
"$LSWS_CTRL" restart >/dev/null 2>&1 || systemctl restart lshttpd

sleep 2

# Verify
if "$LSWS_CTRL" status 2>/dev/null | grep -q running; then
    print_green "✓ OpenLiteSpeed restarted successfully"
else
    print_error "Failed to restart OpenLiteSpeed"
fi

# Final instructions
echo ""
print_green "=========================================="
print_green "Virtual Host Created Successfully!"
print_green "=========================================="
echo ""
echo "Domain:        $DOMAIN"
echo "Document Root: $DOCROOT"
echo "Config File:   $VHOST_CONF"
echo "Logs:          $VHOST_ROOT/logs/"
echo ""
print_info "Next Steps:"
echo "1. Point your DNS A record to this server's IP"
echo "2. Upload your website files to: $DOCROOT"
echo "3. Test PHP: http://$DOMAIN/info.php"
echo ""
print_warn "Security Note: Delete info.php after testing"
echo ""