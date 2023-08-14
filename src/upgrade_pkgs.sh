#!/usr/bin/env bash

grep -E "gotest|golang|github" go.mod|awk '{print "go get "$1}'|sh
