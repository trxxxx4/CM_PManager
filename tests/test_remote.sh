#!/bin/sh
cd "$(dirname "$0")/.." || exit 1

run() {
    echo "=== $1 ==="
    go run ./src --config "tests/$2.toml"
    echo
}

run "Wrong version: error lists real versions" remote_wrong_version
run "Repository URL not found" remote_404
run "Server unreachable" remote_unreachable
