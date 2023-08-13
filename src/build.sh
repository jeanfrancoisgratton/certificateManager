#!/usr/bin/env sh

OUTPUT=/opt/bin

if [ "$#" -gt 0 ]; then
    OUTPUT=$1
fi
go build -o ${OUTPUT}/cm .

