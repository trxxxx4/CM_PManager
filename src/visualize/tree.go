package visualize

import (
	"fmt"
	"io"
	"sort"

	"depgraph/src/graph"
)

// WriteASCIITree writes the dependency graph as a cycle-safe ASCII tree.
func WriteASCIITree(output io.Writer, result *graph.Graph, reverse bool) {
	children := adjacency(result.Edges)
	title := "Dependency tree"
	if reverse {
		title = "Reverse dependency tree"
	}
	fmt.Fprintf(output, "%s:\n%s\n", title, result.Root)
	active := map[string]bool{result.Root: true}
	visited := map[string]bool{result.Root: true}
	writeChildren(output, result.Root, "", children, active, visited)
}

// adjacency groups graph edges by their source vertex.
func adjacency(edges []graph.Edge) map[string][]string {
	children := make(map[string][]string)
	for _, edge := range edges {
		children[edge.From] = append(children[edge.From], edge.To)
	}
	for name := range children {
		sort.Strings(children[name])
	}
	return children
}

// writeChildren prints one tree level and marks repeated or cyclic vertices.
func writeChildren(
	output io.Writer,
	parent, prefix string,
	children map[string][]string,
	active, visited map[string]bool,
) {
	items := children[parent]
	for index, child := range items {
		last := index == len(items)-1
		branch, childPrefix := treeBranch(last, prefix)
		writeChild(output, child, prefix+branch, childPrefix, children, active, visited)
	}
}

// treeBranch returns the branch marker and prefix for the next level.
func treeBranch(last bool, prefix string) (string, string) {
	if last {
		return "`-- ", prefix + "    "
	}
	return "|-- ", prefix + "|   "
}

// writeChild prints a vertex and recursively expands it when possible.
func writeChild(
	output io.Writer,
	name, line, childPrefix string,
	children map[string][]string,
	active, visited map[string]bool,
) {
	fmt.Fprint(output, line, name)
	if active[name] {
		fmt.Fprintln(output, " [cycle]")
		return
	}
	if visited[name] {
		fmt.Fprintln(output, " [already shown]")
		return
	}
	fmt.Fprintln(output)
	active[name], visited[name] = true, true
	writeChildren(output, name, childPrefix, children, active, visited)
	delete(active, name)
}
