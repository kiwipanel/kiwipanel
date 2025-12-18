#!/bin/bash
set -euo pipefail

BASE="/opt/kiwipanel"
BIN="$BASE/bin"
META="$BASE/meta"
WRAPPER="/usr/local/bin/kiwipanel"

EXPECTED_USER="kiwipanel"
EXPECTED_GROUP="kiwisecure"
EXPECTED_MODE="750"

need_root() {
  if [ "$(id -u)" -ne 0 ]; then
    echo "This action requires root privileges" >&2
    exit 1
  fi
}

repair_wrapper() {
  need_root
  if [ -f "$WRAPPER" ]; then
    chown root:root "$WRAPPER"
    chmod 755 "$WRAPPER"
    echo "✔ Wrapper repaired"
  else
    echo "⚠ Wrapper missing, reinstall manually"
  fi
}

repair_binary() {
  need_root
  for slot in A B; do
    f="$BIN/kiwipanel-$slot"
    if [ -f "$f" ]; then
      chown "$EXPECTED_USER:$EXPECTED_GROUP" "$f"
      chmod "$EXPECTED_MODE" "$f"
      echo "✔ Repaired binary $slot"
    fi
  done
  # Ensure symlink points to active slot
  if [ -f "$META/active_slot" ]; then
    active=$(cat "$META/active_slot")
    ln -sfn "$BIN/kiwipanel-$active" "$BIN/kiwipanel"
    echo "✔ Symlink updated to active slot ($active)"
  fi
}

repair_paths() {
  need_root
  chown -R "$EXPECTED_USER:$EXPECTED_GROUP" "$BASE"
  chmod -R 750 "$BASE"
  echo "✔ Base directory permissions repaired"
}

repair_service() {
  need_root
  if systemctl list-unit-files | grep -q kiwipanel.service; then
    systemctl daemon-reexec
    systemctl daemon-reload
    systemctl enable kiwipanel.service
    systemctl restart kiwipanel.service
    echo "✔ systemd service repaired"
  else
    echo "⚠ systemd unit missing, create manually"
  fi
}

main() {
  echo "KiwiPanel Repair"
  echo "----------------"

  repair_wrapper
  repair_binary
  repair_paths
  repair_service

  echo "✔ Repair completed"
}

main
