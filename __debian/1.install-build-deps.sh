#!/usr/bin/env bash

# NOTE:
# -----
# This docker container has been stripped down as much as possible, and works for GO software
# The following extra packages might be needed for languages other than GO :
# sudo apt install -y g++ fakeroot devscripts build-essential

echo "Installing dependencies";echo
sudo apt-get update && sudo apt update -y
echo;echo;echo "Done. Now installing the Go binaries"
sudo rm -rf /opt/go-versions ; sudo mkdir -p /opt/go-versions
sudo /opt/bin/install_golang.sh `cat ../go.version` amd64
