
#!/bin/bash
echo "=== Remote mode ==="
go run ./src/main.go --config tests/valid_remote.toml
echo
echo "=== Test mode ==="
go run ./src/main.go --config tests/valid_test.toml
echo
echo "=== Default config.toml ==="
go run ./src/main.go