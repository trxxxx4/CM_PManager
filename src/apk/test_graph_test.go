package apk

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTestGraph(t *testing.T) {
	path := filepath.Join(t.TempDir(), "graph.txt")
	content := "# graph\nA: B C\nB: C\nC:\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	packages, err := LoadTestRepository(path, "1.0")
	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 3 || packages[0].Name != "A" || packages[0].Version != "1.0" {
		t.Fatalf("unexpected packages: %#v", packages)
	}
	if len(packages[0].DirectDependencies()) != 2 {
		t.Fatalf("unexpected dependencies: %#v", packages[0].Dependencies)
	}
}

func TestLoadTestGraphRejectsInvalidPackageName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "graph.txt")
	if err := os.WriteFile(path, []byte("lowercase: A\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadTestRepository(path, "1.0"); err == nil {
		t.Fatal("expected invalid graph error")
	}
}

func TestLoadTestGraphCanStartWithPackageP(t *testing.T) {
	path := filepath.Join(t.TempDir(), "graph.txt")
	if err := os.WriteFile(path, []byte("P: A\nA:\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	packages, err := LoadTestRepository(path, "1.0")
	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 2 || packages[0].Name != "P" {
		t.Fatalf("unexpected packages: %#v", packages)
	}
}
