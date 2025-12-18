#!/bin/bash
set -euo pipefail

# =========================
# Configuration
# =========================
BASE="/opt/kiwipanel"
BIN="$BASE/bin"
META="$BASE/meta"
SCRIPT_DIR="$BASE/scripts"

SERVICE="kiwipanel.service"

EXPECTED_USER="kiwipanel"
EXPECTED_GROUP="kiwisecure"
EXPECTED_MODE="750"

DOWNLOAD_URL="https://kiwipanel.com/releases/latest/kiwipanel-linux-amd64"
CHECKSUM_URL="${DOWNLOAD_URL}.sha256"

TMP="/tmp/kiwipanel.new"

# =========================
# Helpers
# =========================
fatal() {
  echo >&2
  echo "KiwiPanel update failed:" >&2
  echo "  $1" >&2
  echo >&2
  exit 1
}

need_root() {
  if [ "$(id -u)" -ne 0 ]; then
    fatal "This command must be run as root"
  fi
}

slot_active() {
  cat "$META/active_slot" 2>/dev/null || echo A
}

slot_inactive() {
  [ "$(slot_active)" = "A" ] && echo B || echo A
}

bin_for_slot() {
  echo "$BIN/kiwipanel-$1"
}

verify_binary() {
  local f="$1"

  [ -f "$f" ] || fatal "Binary missing: $f"

  local o g m
  o=$(stat -c %U "$f")
  g=$(stat -c %G "$f")
  m=$(stat -c %a "$f")

  [ "$o" = "$EXPECTED_USER" ] || fatal "Owner $o != $EXPECTED_USER"
  [ "$g" = "$EXPECTED_GROUP" ] || fatal "Group $g != $EXPECTED_GROUP"
  [ "$m" -ge "$EXPECTED_MODE" ] || fatal "Mode $m < $EXPECTED_MODE"

  file "$f" | grep -q ELF || fatal "Binary is not a valid ELF executable"
}

rollback() {
  echo "Rolling back update..."
  local prev
  prev="$(slot_active)"

  ln -sfn "$(bin_for_slot "$prev")" "$BIN/kiwipanel"
  systemctl restart "$SERVICE" || true

  fatal "Rollback completed"
}

# =========================
# Update logic
# =========================
update_online() {
  need_root

  mkdir -p "$META"

  local slot target
  slot="$(slot_inactive)"
  target="$(bin_for_slot "$slot")"

  echo "Downloading new binary (slot $slot)..."
  curl -fsSL "$DOWNLOAD_URL" -o "$TMP"
  curl -fsSL "$CHECKSUM_URL" -o "$TMP.sha256"

  echo "Verifying checksum..."
  sha256sum -c "$TMP.sha256" || fatal "Checksum verification failed"

  mv "$TMP" "$target"
  chown "$EXPECTED_USER:$EXPECTED_GROUP" "$target"
  chmod "$EXPECTED_MODE" "$target"

  verify_binary "$target"

  echo "Activating slot $slot"
  ln -sfn "$target" "$BIN/kiwipanel"
  echo "$slot" > "$META/active_slot"

  echo "Restarting service..."
  systemctl restart "$SERVICE" || rollback

  echo "✔ KiwiPanel updated successfully (slot $slot)"
}

update_offline() {
  need_root
  local src="$1"

  [ -f "$src" ] || fatal "File not found: $src"

  local slot target
  slot="$(slot_inactive)"
  target="$(bin_for_slot "$slot")"

  cp "$src" "$target"
  chown "$EXPECTED_USER:$EXPECTED_GROUP" "$target"
  chmod "$EXPECTED_MODE" "$target"

  verify_binary "$target"

  ln -sfn "$target" "$BIN/kiwipanel"
  echo "$slot" > "$META/active_slot"

  systemctl restart "$SERVICE" || rollback

  echo "✔ Offline update successful (slot $slot)"
}

# =========================
# Entry
# =========================
case "${1:-}" in
  update)
    update_online
    ;;
  update-offline)
    update_offline "${2:-}"
    ;;
  *)
    fatal "Usage: update.sh {update|update-offline <file>}"
    ;;
esac
