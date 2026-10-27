#!/bin/sh
cd "$(dirname "$0")/.." || exit 1

run() {
    echo "=== $1 ==="
    go run ./src --config "tests/$2.toml"
    echo
}

run "Package not found" err_package
run "Version not found" err_version
