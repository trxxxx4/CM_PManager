#!/bin/bash
echo === Invalid name, URL, version and depth ===
go run ./src/main.go --config tests/invalid_values.toml
echo.
echo === Invalid repository mode ===
go run ./src/main.go --config tests/invalid_mode.toml
echo.
echo === Test repository file not found ===
go run ./src/main.go --config tests/invalid_repo_file.toml
echo.
echo === Test repository is a directory ===
go run ./src/main.go --config tests/invalid_repo_dir.toml
