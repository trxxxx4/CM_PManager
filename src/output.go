package main

import (
	"fmt"
	"io"

	"depgraph/src/graph"
)

// printGraph writes the nodes and directed edges of the dependency graph.
func printGraph(w io.Writer, result *graph.Graph, reverse bool) {
	if reverse {
		printReverseGraph(w, result)
		return
	}
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

// printReverseGraph shows traversal edges as "target <- dependent".
func printReverseGraph(w io.Writer, result *graph.Graph) {
	fmt.Fprintf(w, "Reverse dependency graph: %s\n", result.Root)
	fmt.Fprintf(w, "Packages (%d):\n", len(result.Nodes))
	for _, name := range result.Nodes {
		fmt.Fprintf(w, "  %s\n", name)
	}
	fmt.Fprintf(w, "Reverse dependencies (%d):\n", len(result.Edges))
	if len(result.Edges) == 0 {
		fmt.Fprintln(w, "  none")
	}
	for _, edge := range result.Edges {
		fmt.Fprintf(w, "  %s <- %s\n", edge.From, edge.To)
	}
}
