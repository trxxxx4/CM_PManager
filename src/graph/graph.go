// Package graph builds dependency graphs from APK package records.
package graph

import (
	"fmt"
	"sort"

	"depgraph/src/apk"
)

// Edge describes a directed dependency from one package to another.
type Edge struct {
	From string
	To   string
}

// Graph contains all packages and dependency edges reached by DFS.
type Graph struct {
	Root  string
	Nodes []string
	Edges []Edge
}

// builder stores DFS state for one graph construction.
type builder struct {
	packages map[string]*apk.Package
	nodes    map[string]bool
	path     map[string]bool
	expanded map[string]int
	result   *Graph
	maxDepth int
}

// Build traverses package dependencies with recursive depth-first search.
func Build(packages []apk.Package, root string, maxDepth int) (*Graph, error) {
	byName := make(map[string]*apk.Package, len(packages))
	for i := range packages {
		if _, exists := byName[packages[i].Name]; !exists {
			byName[packages[i].Name] = &packages[i]
		}
	}
	if byName[root] == nil {
		return nil, fmt.Errorf("package not found in repository: %s", root)
	}

	state := &builder{
		packages: byName,
		nodes:    make(map[string]bool),
		path:     make(map[string]bool),
		expanded: make(map[string]int),
		result:   &Graph{Root: root},
		maxDepth: maxDepth,
	}
	state.visit(root, 0)
	sort.Strings(state.result.Nodes)
	sort.Slice(state.result.Edges, func(i, j int) bool {
		if state.result.Edges[i].From == state.result.Edges[j].From {
			return state.result.Edges[i].To < state.result.Edges[j].To
		}
		return state.result.Edges[i].From < state.result.Edges[j].From
	})
	return state.result, nil
}

// visit adds reachable dependencies and prevents recursive cycles.
func (b *builder) visit(name string, depth int) {
	if !b.nodes[name] {
		b.nodes[name] = true
		b.result.Nodes = append(b.result.Nodes, name)
	}
	if depth >= b.maxDepth || b.path[name] ||
		b.expanded[name] >= b.maxDepth-depth {
		return
	}

	b.path[name] = true
	b.expanded[name] = b.maxDepth - depth
	defer delete(b.path, name)

	pkg := b.packages[name]
	if pkg == nil {
		return
	}
	for _, dependency := range pkg.DirectDependencies() {
		depName := dependency.Name
		addEdge(b.result, name, depName)
		b.visit(depName, depth+1)
	}
}

// addEdge keeps each directed edge only once.
func addEdge(result *Graph, from, to string) {
	for _, edge := range result.Edges {
		if edge.From == from && edge.To == to {
			return
		}
	}
	result.Edges = append(result.Edges, Edge{From: from, To: to})
}
