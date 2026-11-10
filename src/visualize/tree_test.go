package visualize

import (
	"bytes"
	"strings"
	"testing"

	"depgraph/src/graph"
)

func TestWriteASCIITreeShowsBranchesAndSharedNodes(t *testing.T) {
	result := &graph.Graph{
		Root: "A",
		Edges: []graph.Edge{
			{From: "A", To: "B"},
			{From: "A", To: "C"},
			{From: "B", To: "D"},
			{From: "C", To: "D"},
		},
	}
	var output bytes.Buffer

	WriteASCIITree(&output, result, false)
	for _, part := range []string{
		"Dependency tree:\nA\n",
		"|-- B\n|   `-- D\n",
		"`-- C\n    `-- D [already shown]",
	} {
		if !strings.Contains(output.String(), part) {
			t.Errorf("tree output does not contain %q: %s", part, output.String())
		}
	}
}

func TestWriteASCIITreeMarksCycles(t *testing.T) {
	result := &graph.Graph{
		Root:  "A",
		Edges: []graph.Edge{{From: "A", To: "B"}, {From: "B", To: "A"}},
	}
	var output bytes.Buffer

	WriteASCIITree(&output, result, false)
	if !strings.Contains(output.String(), "`-- A [cycle]") {
		t.Fatalf("cycle marker is missing: %s", output.String())
	}
}

func TestWriteASCIITreeUsesReverseTitle(t *testing.T) {
	result := &graph.Graph{Root: "A"}
	var output bytes.Buffer
	WriteASCIITree(&output, result, true)
	if !strings.Contains(output.String(), "Reverse dependency tree:") {
		t.Fatalf("reverse title is missing: %s", output.String())
	}
}
