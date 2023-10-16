#!/usr/bin/env bash

if [ $# -ne 2 ]; then
    echo "Usage: ./upgradeCM.sh goversion_number architecture"
    echo "Where architecture is either amd64 or arm64"
    exit
fi

# Fetch GO binaries
wget -q https://go.dev/dl/go$1.linux-$2.tar.gz -O /tmp/go.tar.gz && cd /opt && sudo tar zxf /tmp/go.tar.gz && rm /tmp/go.tar.gz

# Clone certificateManager git repo
rm -rf /home/builder/certificateManager 
git clone ssh://git@git.famillegratton.net:9722/devops/certificateManager.git /home/builder/certificateManager
cd /home/builder/certificateManager/src && ./build.sh

# Cleanup
sudo rm -rf /home/builder/certificateManager /opt/go