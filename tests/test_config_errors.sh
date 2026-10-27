#!/bin/sh
cd "$(dirname "$0")/.." || exit 1

run() {
    echo "=== $1 ==="
    go run ./src --config "tests/$2.toml"
    echo
}

run "Invalid name, URL, version and depth" invalid_values
run "Invalid repository mode" invalid_mode
run "Test repository file not found" invalid_repo_file
run "Test repository is a directory" invalid_repo_dir
run "No parameters in the file" empty
run "Config file not found" no_such_config
run "max_depth is a string" invalid_type_depth
run "ascii_tree is a string" invalid_type_ascii
run "package_name is a number" invalid_type_name
run "Broken TOML syntax" broken_syntax
run "Unknown parameter" unknown_key

echo "=== Unexpected argument ==="
go run ./src tests/empty.toml
