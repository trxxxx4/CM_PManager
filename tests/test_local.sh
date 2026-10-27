#!/bin/sh
cd "$(dirname "$0")/.." || exit 1

run() {
    echo "=== $1 ==="
    go run ./src --config "tests/$2.toml"
    echo
}

run "Package with dependencies" local_curl
run "Version constraints" local_openssl
run "Package without dependencies" local_musl
run "Conflicts are not dependencies" local_conflicts
