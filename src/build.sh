#!/usr/bin/env sh

GOROOT=/opt/go
OUTPUT=/opt/bin

if [ "$#" -gt 0 ]; then
    OUTPUT=$1
fi
sudo ${GOROOT}/bin/go build -o ${OUTPUT}/cm .
sudo strip ${OUTPUT}/cm
