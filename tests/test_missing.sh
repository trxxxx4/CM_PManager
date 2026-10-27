#!/bin/bash
echo === No parameters in the file ===
go run ./src/main.go --config tests/empty.toml
echo.
echo === Config file not found ===
go run ./src/main.go --config tests/no_such_config.toml
echo.
echo === Unexpected argument ===
go run ./src/main.go tests/valid_remote.toml
