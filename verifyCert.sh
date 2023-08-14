#!/usr/bin/env bash

openssl x509 -noout -text -in $1.crt | more
