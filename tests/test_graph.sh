#!/bin/sh
set -eu
cd "$(dirname "$0")/.."

run() {
    printf '\n=== %s ===\n' "$1"
    go run ./src --config "tests/$2.toml"
}

run "Cycle handling" graph_cycle
run "Maximum depth" graph_depth
run "Package without dependencies" graph_leaf
