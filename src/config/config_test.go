package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const validConfig = `
package_name = "left-pad"
repository = "https://registry.npmjs.org"
repository_mode = "remote"
version = "1.3.0"
ascii_tree = true
max_depth = 3
`

// writeFile creates a file with the content in a temporary directory.
func writeFile(t *testing.T, name, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	return path
}

// TestLoadValid checks that a correct file is loaded into the struct.
func TestLoadValid(t *testing.T) {
	cfg, err := Load(writeFile(t, "config.toml", validConfig))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Config{
		PackageName:    "left-pad",
		Repository:     "https://registry.npmjs.org",
		RepositoryMode: ModeRemote,
		Version:        "1.3.0",
		ASCIITree:      true,
		MaxDepth:       3,
	}

	if *cfg != want {
		t.Errorf("got %+v, want %+v", *cfg, want)
	}
}

// withDepth returns the valid configuration with another max_depth.
func withDepth(value string) string {
	return strings.Replace(
		validConfig,
		"max_depth = 3",
		"max_depth = "+value,
		1,
	)
}

// TestLoadErrors checks the messages for broken configuration files.
func TestLoadErrors(t *testing.T) {
	unknown := validConfig + "extra = 1\n"
	cases := []struct {
		name    string
		content string
		want    string
	}{
		{"empty file", "", "package_name: required parameter is missing"},
		{"unknown key", unknown, "extra: unknown parameter"},
		{"bad syntax", "package_name = ", "invalid config file"},
		{"wrong type", withDepth(`"x"`), "max_depth"},
		{"bad depth", withDepth("0"), "must be at least"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Load(writeFile(t, "config.toml", tc.content))

			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Errorf("got %v, want error containing %q", err, tc.want)
			}
		})
	}
}

// TestLoadMissingFile checks the error for a nonexistent config file.
func TestLoadMissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "absent.toml"))

	if err == nil || !strings.Contains(err.Error(), "config file not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestValidate checks the exported validation of already built values.
func TestValidate(t *testing.T) {
	cfg := &Config{
		PackageName:    "left-pad",
		Repository:     "https://registry.npmjs.org",
		RepositoryMode: ModeRemote,
		Version:        "1.3.0",
		MaxDepth:       1,
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	cfg.MaxDepth = maxDepthLimit + 1

	if err := cfg.Validate(); err == nil {
		t.Error("expected an error for a too large depth")
	}
}
