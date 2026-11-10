#!/bin/sh
set -eu
cd "$(dirname "$0")/.."

run() {
    printf '\n=== %s ===\n' "$1"
    go run ./src --config "tests/$2.toml" --reverse
}

run "Direct reverse dependencies (depth 1)" reverse_depth
run "Transitive reverse dependencies and cycle" reverse_cycle
run "Package without reverse dependencies" reverse_leaf
