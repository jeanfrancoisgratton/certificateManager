#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

rm -rf src pkg


PKGDEST=/data makepkg --force --syncdeps --noconfirm

for d in src pkg; do
  if [[ -d "$d" ]]; then
    chmod -R u+w "$d" 2>/dev/null || true
    rm -rf "$d"
  fi
done
