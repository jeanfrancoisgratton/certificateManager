#!/usr/bin/env bash
set -euo pipefail

echo "Installing Arch package build dependencies"
sudo pacman -Syu --needed base-devel

echo "Completed"
