#!/usr/bin/env bash

echo "Installing dependencies";echo
sudo apt-get update && sudo apt update -y
echo;echo;echo "Done. Now installing the Go binaries"
sudo rm -rf /opt/go-versions && sudo mkdir -p /opt/go-versions
sudo /opt/bin/install_golang.sh `cat ../go.version` amd64
