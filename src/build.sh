#!/usr/bin/env sh

GOROOT=/opt/go
OUTPUT=/opt/bin
BINARY=cm

# Get git branch's name, replace / with _
BRANCH=`git rev-parse --abbrev-ref HEAD`
BRANCH=$(echo "$BRANCH" | tr '/' '_')

if [ "$BRANCH" = "master" ] || [ "$BRANCH" = "main" ] || [ "$BRANCH" = "develop" ]; then
    FULLNAME="$BINARY"
else
    FULLNAME="$BINARY-$BRANCH"
fi

if [ "$#" -gt 0 ]; then
    OUTPUT=$1
fi

go build -o $OUTPUT/$FULLNAME .
