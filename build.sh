#!/bin/bash
set -e

for dir in cmd/*/; do
    dir=${dir%/}
    if grep -q '^package main$' $dir/*.go 2>/dev/null; then
        echo "building $dir"
        CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o $dir/$(basename $dir) ./$dir
        docker build -t antihax/evedata-$(basename $dir) -f Dockerfile.$(basename $dir) .
    else
        echo "(skipped $dir)"
    fi
done
