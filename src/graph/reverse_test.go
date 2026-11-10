package graph

import (
	"fmt"
	"reflect"
	"testing"

	"depgraph/src/apk"
)

// TestBuildReverseCases covers branches, depth, cycles and unrelated packages.
func TestBuildReverseCases(t *testing.T) {
	packages := []apk.Package{
		{Name: "A", Dependencies: []apk.Dependency{{Name: "B"}, {Name: "C"}}},
		{Name: "B", Dependencies: []apk.Dependency{{Name: "D"}}},
		{Name: "C", Dependencies: []apk.Dependency{{Name: "D"}}},
		{Name: "D", Dependencies: []apk.Dependency{{Name: "A"}}},
		{Name: "E"},
	}
	cases := []struct {
		name  string
		depth int
		nodes []string
		edges []Edge
	}{
		{"D", 1, []string{"B", "C", "D"}, []Edge{{"D", "B"}, {"D", "C"}}},
		{"D", 2, []string{"A", "B", "C", "D"}, []Edge{
			{"B", "A"}, {"C", "A"}, {"D", "B"}, {"D", "C"},
		}},
		{"D", 10, []string{"A", "B", "C", "D"}, []Edge{
			{"A", "D"}, {"B", "A"}, {"C", "A"}, {"D", "B"}, {"D", "C"},
		}},
		{"E", 10, []string{"E"}, nil},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s-depth-%d", tc.name, tc.depth), func(t *testing.T) {
			result, err := BuildReverse(packages, tc.name, tc.depth)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result.Nodes, tc.nodes) || !equalEdges(result.Edges, tc.edges) {
				t.Fatalf("got %+v, want nodes %v edges %v", result, tc.nodes, tc.edges)
			}
		})
	}
}

// equalEdges treats an empty slice and nil as the same empty edge set.
func equalEdges(got, want []Edge) bool {
	return len(got) == 0 && len(want) == 0 || reflect.DeepEqual(got, want)
}

// TestBuildReverseExcludesConflicts checks APK conflict handling.
func TestBuildReverseExcludesConflicts(t *testing.T) {
	packages := []apk.Package{
		{Name: "A"},
		{Name: "B", Dependencies: []apk.Dependency{{Name: "A", Conflict: true}}},
		{Name: "C", Dependencies: []apk.Dependency{{Name: "A"}, {Name: "A"}}},
	}
	result, err := BuildReverse(packages, "A", 3)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result.Nodes, []string{"A", "C"}) || len(result.Edges) != 1 {
		t.Fatalf("unexpected graph: %+v", result)
	}
}

// TestBuildReverseUnknownRoot checks the error for a missing package.
func TestBuildReverseUnknownRoot(t *testing.T) {
	if _, err := BuildReverse(nil, "Z", 2); err == nil {
		t.Fatal("expected missing package error")
	}
}
