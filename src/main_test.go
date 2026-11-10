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

// TestRunPrintsDependencies checks the output for packages of the index.
func TestRunPrintsDependencies(t *testing.T) {
	cases := []struct {
		name    string
		version string
		want    string
	}{
		{
			"curl", "8.5.0-r0",
			"Dependency graph: curl\n" +
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
			"musl", "1.2.4-r4",
			"Dependency graph: musl\nPackages (1):\n  musl\nDependencies (0):\n  none\n",
		},
	}

	for _, tc := range cases {
		var stdout, stderr bytes.Buffer
		config := writeConfig(t, tc.name, tc.version)

		code := run([]string{"--config", config}, &stdout, &stderr)

		if code != exitOK || stdout.String() != tc.want {
			t.Errorf("%s: code %d, output %q, stderr %q",
				tc.name, code, stdout.String(), stderr.String())
		}
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
	var stdout, stderr bytes.Buffer
	code := run([]string{"--config", "../tests/reverse_depth.toml", "--reverse"}, &stdout, &stderr)
	want := "Reverse dependency graph: D\nPackages (3):\n" +
		"  B\n  C\n  D\nReverse dependencies (2):\n" +
		"  D <- B\n  D <- C\n"
	if code != exitOK || stdout.String() != want {
		t.Fatalf("code %d, output %q, stderr %q", code, stdout.String(), stderr.String())
	}
}
