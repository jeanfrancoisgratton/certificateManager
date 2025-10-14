#!/usr/bin/env bash

# NOTE:
# -----
echo "# This docker container has been stripped down as much as possible, and works for GO software"
echo "# The following extra packages might be needed for languages other than GO :"
echo "# sudo apt install -y g++ fakeroot devscripts build-essential";echo

echo "Installing dependencies";echo
sudo apt-get update && sudo apt update -y
echo;echo;echo "Done. Now installing the Go binaries"
sudo /opt/bin/install_golang.sh 1.25.2 amd64
