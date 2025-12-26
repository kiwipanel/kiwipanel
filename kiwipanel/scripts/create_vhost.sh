#!/bin/bash
set -e

USER="$1"
DOMAIN="$2"

DOCROOT="/home/$USER/$DOMAIN/public_html"
VHOST_CONF="/opt/kiwipanel/config/lsws/vhosts/$DOMAIN.conf"

if ! id "$USER" &>/dev/null; then
    echo "User $USER does not exist"
    exit 1
fi

mkdir -p "$DOCROOT"
chown -R "$USER:$USER" "/home/$USER/$DOMAIN"
chmod 750 "/home/$USER"

cat > "$VHOST_CONF" <<EOF
vhRoot                  $DOCROOT
configFile              \$VH_ROOT/conf/vhconf.conf
allowSymbolLink         0
enableScript            1
restrained              1
context / {
    docRoot             $DOCROOT
}
EOF

chown kiwipanel:kiwisecure "$VHOST_CONF"
chmod 640 "$VHOST_CONF"

echo "VHost created for $DOMAIN"
