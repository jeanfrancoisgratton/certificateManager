#!/usr/bin/env sh

GOROOT=/opt/go
OUTPUT=/opt/bin

arch=$(uname -m)


if [ "$#" -gt 0 ]; then
    OUTPUT=$1
fi

sudo ${GOROOT}/bin/go build -o ${OUTPUT}/cm .

# Somehow fails on ARM64
if [ "$arch" = "x86_64" ]; then
    sudo strip ${OUTPUT}/cm
fi
