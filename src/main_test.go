package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testIndex = `P:curl
V:8.5.0-r0
D:ca-certificates so:libc.musl-x86_64.so.1 !curl-old libcurl>=8.5.0

P:musl
V:1.2.4-r4
`

// writeConfig creates an index and a config that points to it.
func writeConfig(t *testing.T, name, version string) string {
	t.Helper()

	dir := t.TempDir()
	t.Setenv("DEPGRAPH_NO_OPEN", "1")
	t.Cleanup(cleanupGraphArtifacts)
	index := filepath.Join(dir, "APKINDEX")
	config := filepath.Join(dir, "config.toml")

	if err := os.WriteFile(index, []byte(testIndex), 0o600); err != nil {
		t.Fatal(err)
	}

	content := "package_name = '" + name + "'\n" +
		"repository = '" + index + "'\n" +
		"repository_mode = 'test'\n" +
		"version = '" + version + "'\n" +
		"ascii_tree = true\n" +
		"max_depth = 3\n"

	if err := os.WriteFile(config, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	return config
}

// outputCase describes one package output expectation.
type outputCase struct {
	name    string
	version string
	want    string
}

// TestRunPrintsDependencies checks output for packages in the test index.
func TestRunPrintsDependencies(t *testing.T) {
	for _, tc := range dependencyOutputCases() {
		t.Run(tc.name, func(t *testing.T) {
			checkDependencyOutput(t, tc)
		})
	}
}

// dependencyOutputCases returns expected output prefixes for sample packages.
func dependencyOutputCases() []outputCase {
	return []outputCase{
		{
			name: "curl", version: "8.5.0-r0",
			want: "Dependency graph: curl\n" +
				"Packages (4):\n" +
				"  ca-certificates\n" +
				"  curl\n" +
				"  libcurl\n" +
				"  so:libc.musl-x86_64.so.1\n" +
				"Dependencies (3):\n" +
				"  curl -> ca-certificates\n" +
				"  curl -> libcurl\n" +
				"  curl -> so:libc.musl-x86_64.so.1\n",
		},
		{
			name: "musl", version: "1.2.4-r4",
			want: "Dependency graph: musl\nPackages (1):\n  musl\n" +
				"Dependencies (0):\n  none\n",
		},
	}
}

// checkDependencyOutput verifies graph, DOT and ASCII-tree output.
func checkDependencyOutput(t *testing.T, tc outputCase) {
	var stdout, stderr bytes.Buffer
	configPath := writeConfig(t, tc.name, tc.version)
	code := run([]string{"--config", configPath}, &stdout, &stderr)
	output := stdout.String()
	if code != exitOK || !strings.HasPrefix(output, tc.want) {
		t.Fatalf("code %d, output %q, stderr %q", code, output, stderr.String())
	}
	checkOutputContains(t, output, "Graphviz DOT: dependency_graph_"+tc.name+".dot")
	checkOutputContains(t, output, "Dependency tree:")
}

// checkOutputContains fails when expected text is absent from command output.
func checkOutputContains(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.Contains(output, expected) {
		t.Errorf("output does not contain %q: %s", expected, output)
	}
}

// TestRunErrors checks the exit codes of the failures.
func TestRunErrors(t *testing.T) {
	extra := []string{"extra"}
	absent := []string{"--config", "absent.toml"}
	cases := []struct {
		name string
		args []string
		code int
		want string
	}{
		{"unknown package", nil, exitDataError, "package not found"},
		{"unexpected argument", extra, exitUsageError, "unexpected"},
		{"no config", absent, exitConfigError, "not found"},
	}

	for _, tc := range cases {
		var stdout, stderr bytes.Buffer
		args := tc.args

		if args == nil {
			args = []string{"--config", writeConfig(t, "absent", "1.0")}
		}

		code := run(args, &stdout, &stderr)

		if code != tc.code || !strings.Contains(stderr.String(), tc.want) {
			t.Errorf("%s: code %d, stderr %q", tc.name, code, stderr.String())
		}
	}
}

// TestRunReverse checks CLI output for direct reverse dependencies.
func TestRunReverse(t *testing.T) {
	t.Setenv("DEPGRAPH_NO_OPEN", "1")
	t.Cleanup(cleanupGraphArtifacts)
	var stdout, stderr bytes.Buffer
	code := run([]string{"--config", "../tests/reverse_depth.toml", "--reverse"}, &stdout, &stderr)
	want := "Reverse dependency graph: D\nPackages (3):\n" +
		"  B\n  C\n  D\nReverse dependencies (2):\n" +
		"  D <- B\n  D <- C\n"
	if code != exitOK || !strings.HasPrefix(stdout.String(), want) {
		t.Fatalf("code %d, output %q, stderr %q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "reverse_dependency_graph_D.png") {
		t.Fatalf("reverse graph image path is missing: %s", stdout.String())
	}
}

// cleanupGraphArtifacts removes output files created by command tests.
func cleanupGraphArtifacts() {
	matches, _ := filepath.Glob("*dependency_graph_*")
	for _, path := range matches {
		_ = os.Remove(path)
	}
}
