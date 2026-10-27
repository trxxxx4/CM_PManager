#!/bin/bash
echo === max_depth is a string ===
go run ./src/main.go --config tests/invalid_type_depth.toml
echo.
echo === ascii_tree is a string ===
go run ./src/main.go --config tests/invalid_type_ascii.toml
echo.
echo === package_name is a number ===
go run ./src/main.go --config tests/invalid_type_name.toml
