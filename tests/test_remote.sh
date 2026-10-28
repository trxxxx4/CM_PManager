#!/bin/sh
cd "$(dirname "$0")/.." || exit 1

run() {
    echo "=== $1 ==="
    go run ./src --config "tests/$2.toml"
    echo
}

run "Real repository: busybox" remote_busybox
run "Real repository: busybox-suid" remote_busybox_suid
