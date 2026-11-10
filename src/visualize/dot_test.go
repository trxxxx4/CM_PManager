package visualize

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"depgraph/src/graph"
)

func TestToDOTWritesNodesAndEdges(t *testing.T) {
	result := &graph.Graph{
		Root:  "A",
		Nodes: []string{"A", "B", "C"},
		Edges: []graph.Edge{{From: "A", To: "B"}, {From: "B", To: "C"}},
	}

	dot := ToDOT(result, false)
	for _, part := range []string{
		"digraph dependencies {",
		`label="Dependencies of A"`,
		`"A" [label="A"]`,
		`"A" -> "B";`,
		`"B" -> "C";`,
	} {
		if !strings.Contains(dot, part) {
			t.Errorf("DOT output does not contain %q: %s", part, dot)
		}
	}
}

func TestToDOTLabelsReverseGraph(t *testing.T) {
	result := &graph.Graph{Root: "A", Nodes: []string{"A"}}
	if !strings.Contains(ToDOT(result, true), "Reverse dependencies of A") {
		t.Fatal("reverse graph title is missing")
	}
}

const emptyFileSize int64 = 0

func TestGenerateFilesCreatesDOTAndPNG(t *testing.T) {
	if _, err := exec.LookPath("dot"); err != nil {
		t.Skip("Graphviz is not installed")
	}
	dir := t.TempDir()
	dotPath := filepath.Join(dir, "graph.dot")
	imagePath := filepath.Join(dir, "graph.png")
	source := "digraph sample { A -> B }"

	if err := GenerateFiles(dotPath, imagePath, source); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{dotPath, imagePath} {
		info, err := os.Stat(path)
		if err != nil || info.Size() == emptyFileSize {
			t.Errorf("artifact %s was not created", path)
		}
	}
}

func TestFileNamesSanitizePackageName(t *testing.T) {
	dotPath, imagePath := FileNames("so:libc.so.6", false)
	if dotPath != "dependency_graph_so_libc_so_6.dot" {
		t.Fatalf("unexpected DOT path: %s", dotPath)
	}
	if imagePath != "dependency_graph_so_libc_so_6.png" {
		t.Fatalf("unexpected image path: %s", imagePath)
	}
}
