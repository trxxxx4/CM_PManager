package main

import (
	"fmt"
	"io"

	"depgraph/src/graph"
)

// printGraph writes the nodes and directed edges of the dependency graph.
func printGraph(w io.Writer, result *graph.Graph) {
	fmt.Fprintf(w, "Dependency graph: %s\n", result.Root)
	fmt.Fprintf(w, "Packages (%d):\n", len(result.Nodes))
	for _, name := range result.Nodes {
		fmt.Fprintf(w, "  %s\n", name)
	}
	fmt.Fprintf(w, "Dependencies (%d):\n", len(result.Edges))
	if len(result.Edges) == 0 {
		fmt.Fprintln(w, "  none")
		return
	}
	for _, edge := range result.Edges {
		fmt.Fprintf(w, "  %s -> %s\n", edge.From, edge.To)
	}
}
