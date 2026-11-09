package graph

import (
	"reflect"
	"testing"

	"depgraph/src/apk"
)

func TestBuildTraversesDependencies(t *testing.T) {
	packages := []apk.Package{
		{Name: "A", Dependencies: []apk.Dependency{{Name: "B"}, {Name: "C"}}},
		{Name: "B", Dependencies: []apk.Dependency{{Name: "D"}}},
		{Name: "C", Dependencies: []apk.Dependency{{Name: "D"}}},
		{Name: "D"},
	}

	result, err := Build(packages, "A", 3)
	if err != nil {
		t.Fatal(err)
	}

	wantNodes := []string{"A", "B", "C", "D"}
	if !reflect.DeepEqual(result.Nodes, wantNodes) {
		t.Fatalf("nodes = %v, want %v", result.Nodes, wantNodes)
	}
	if len(result.Edges) != 4 {
		t.Fatalf("got %d edges, want 4", len(result.Edges))
	}
}

func TestBuildRespectsMaximumDepth(t *testing.T) {
	packages := []apk.Package{
		{Name: "A", Dependencies: []apk.Dependency{{Name: "B"}}},
		{Name: "B", Dependencies: []apk.Dependency{{Name: "C"}}},
		{Name: "C", Dependencies: []apk.Dependency{{Name: "D"}}},
		{Name: "D"},
	}

	result, err := Build(packages, "A", 2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result.Nodes, []string{"A", "B", "C"}) {
		t.Fatalf("unexpected nodes: %v", result.Nodes)
	}
	if len(result.Edges) != 2 {
		t.Fatalf("got %d edges, want 2", len(result.Edges))
	}
}

func TestBuildStopsCycles(t *testing.T) {
	packages := []apk.Package{
		{Name: "A", Dependencies: []apk.Dependency{{Name: "B"}}},
		{Name: "B", Dependencies: []apk.Dependency{{Name: "C"}}},
		{Name: "C", Dependencies: []apk.Dependency{{Name: "A"}}},
	}

	result, err := Build(packages, "A", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Nodes) != 3 || len(result.Edges) != 3 {
		t.Fatalf("unexpected graph: nodes=%v edges=%v", result.Nodes, result.Edges)
	}
}

func TestBuildRejectsUnknownRoot(t *testing.T) {
	_, err := Build(nil, "A", 1)
	if err == nil {
		t.Fatal("expected an error for unknown root")
	}
}
