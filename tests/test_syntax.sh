#!/bin/bash
echo === Broken TOML syntax ===
go run ./src/main.go --config tests/broken_syntax.toml
echo.
echo === Unknown parameter ===
go run ./src/main.go --config tests/unknown_key.toml
