#!/bin/sh
set -eu

PROJECT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$PROJECT_DIR"

if ! command -v dot >/dev/null 2>&1; then
    echo "Graphviz is required. Install it before running this test." >&2
    exit 1
fi

for package in a b e; do
    DEPGRAPH_NO_OPEN=1 ./run.sh --config "tests/visualize_${package}.toml"
done

for package in A B E; do
    test -s "dependency_graph_${package}.dot"
    test -s "dependency_graph_${package}.png"
done

echo "Visualization scenarios passed for packages A, B and E."
