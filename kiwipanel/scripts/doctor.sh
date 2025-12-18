#!/bin/bash
set -euo pipefail

BASE="/opt/kiwipanel"
BIN="$BASE/bin"
META="$BASE/meta"
WRAPPER="/usr/local/bin/kiwipanel"

EXPECTED_USER="kiwipanel"
EXPECTED_GROUP="kiwisecure"
EXPECTED_MODE="750"

ok()   { echo "✔ $1"; }
warn() { echo "⚠ $1"; }
fail() { echo "✘ $1"; }

check_wrapper() {
  [ -f "$WRAPPER" ] || { fail "CLI wrapper missing"; return 1; }
  [ "$(stat -c %U "$WRAPPER")" = "root" ] || fail "Wrapper not owned by root"
  [ "$(stat -c %a "$WRAPPER")" -eq 755 ] || warn "Wrapper permissions not 755"
  ok "CLI wrapper"
}

check_binary() {
  local f="$BIN/kiwipanel"

  [ -L "$f" ] || { fail "Binary symlink missing"; return 1; }

  local target
  target="$(readlink -f "$f")"

  [ -f "$target" ] || { fail "Binary target missing"; return 1; }

  file "$target" | grep -q ELF || { fail "Binary is not ELF"; return 1; }

  local o g m
  o=$(stat -c %U "$target")
  g=$(stat -c %G "$target")
  m=$(stat -c %a "$target")

  [ "$o" = "$EXPECTED_USER" ] || fail "Binary owner invalid ($o)"
  [ "$g" = "$EXPECTED_GROUP" ] || fail "Binary group invalid ($g)"
  [ "$m" -ge "$EXPECTED_MODE" ] || fail "Binary permissions too weak ($m)"

  ok "Binary & permissions"
}

check_slots() {
  [ -f "$META/active_slot" ] || { fail "active_slot missing"; return 1; }

  local slot
  slot="$(cat "$META/active_slot")"

  [ "$slot" = "A" ] || [ "$slot" = "B" ] || {
    fail "Invalid active slot value: $slot"
    return 1
  }

  [ -f "$BIN/kiwipanel-$slot" ] || fail "Active slot binary missing"

  ok "A/B slot configuration"
}

check_paths() {
  for d in "$BASE" "$BIN" "$META"; do
    [ -d "$d" ] || fail "Missing directory: $d"
    [ -w "$d" ] || warn "Directory not writable: $d"
  done
  ok "Paths"
}

check_service() {
  if systemctl list-unit-files | grep -q kiwipanel.service; then
    ok "systemd unit exists"
  else
    fail "systemd unit missing"
  fi

  systemctl is-active --quiet kiwipanel && ok "Service running" || warn "Service not running"
}

main() {
  echo "KiwiPanel Doctor"
  echo "----------------"

  failures=0

  check_wrapper || failures=$((failures+1))
  check_binary  || failures=$((failures+1))
  check_slots   || failures=$((failures+1))
  check_paths   || failures=$((failures+1))
  check_service || failures=$((failures+1))

  echo
  if [ "$failures" -eq 0 ]; then
    ok "System healthy"
    exit 0
  else
    fail "$failures problem(s) detected"
    exit 1
  fi
}

main
