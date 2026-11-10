// Package visualize formats dependency graphs for Graphviz and terminals.
package visualize

import (
	"fmt"
	"strconv"
	"strings"

	"depgraph/src/graph"
)

// ToDOT converts a dependency graph into Graphviz DOT notation.
func ToDOT(result *graph.Graph, reverse bool) string {
	var output strings.Builder
	writeGraphHeader(&output, result.Root, reverse)
	writeNodes(&output, result.Nodes)
	writeEdges(&output, result.Edges)
	output.WriteString("}\n")
	return output.String()
}

// writeGraphHeader writes graph attributes and its descriptive title.
func writeGraphHeader(output *strings.Builder, root string, reverse bool) {
	title := "Dependencies of " + root
	if reverse {
		title = "Reverse dependencies of " + root
	}
	fmt.Fprintf(output, "digraph dependencies {\n")
	fmt.Fprintf(output, "  graph [rankdir=LR, label=%s, labelloc=t];\n", strconv.Quote(title))
	output.WriteString("  node [shape=box, style=rounded];\n")
}

// writeNodes writes every graph vertex with a safe quoted identifier.
func writeNodes(output *strings.Builder, nodes []string) {
	for _, name := range nodes {
		quoted := strconv.Quote(name)
		fmt.Fprintf(output, "  %s [label=%s];\n", quoted, quoted)
	}
}

// writeEdges writes directed graph edges.
func writeEdges(output *strings.Builder, edges []graph.Edge) {
	for _, edge := range edges {
		fmt.Fprintf(
			output,
			"  %s -> %s;\n",
			strconv.Quote(edge.From),
			strconv.Quote(edge.To),
		)
	}
}
