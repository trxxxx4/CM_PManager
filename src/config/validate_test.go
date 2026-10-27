package config

import (
	"path/filepath"
	"strings"
	"testing"
)

type valueCase struct {
	name    string
	value   string
	wantErr bool
}

// runValueCases runs check on every value and compares the outcome.
func runValueCases(t *testing.T, cases []valueCase, run func(string) error) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := run(tc.value)

			if (err != nil) != tc.wantErr {
				t.Errorf("value %q: got error %v, want error: %t",
					tc.value, err, tc.wantErr)
			}
		})
	}
}

// TestCheckPackageName covers valid and invalid package names.
func TestCheckPackageName(t *testing.T) {
	cases := []valueCase{
		{"plain", "left-pad", false},
		{"scoped", "@scope/pkg", false},
		{"empty", "", true},
		{"blank", "   ", true},
		{"space inside", "bad name", true},
		{"bad symbol", "pkg!", true},
		{"too long", strings.Repeat("a", maxPackageNameLength+1), true},
	}

	runValueCases(t, cases, func(v string) error {
		return checkPackageName(&Config{PackageName: v})
	})
}

// TestCheckVersion covers valid and invalid versions.
func TestCheckVersion(t *testing.T) {
	cases := []valueCase{
		{"semver", "1.3.0", false},
		{"prerelease", "2.0.0-beta.1", false},
		{"empty", "", true},
		{"space", "1.0 0", true},
		{"leading symbol", "-1", true},
		{"too long", strings.Repeat("1", maxVersionLength+1), true},
	}

	runValueCases(t, cases, func(v string) error {
		return checkVersion(&Config{Version: v})
	})
}

// TestCheckMode covers the known and unknown repository modes.
func TestCheckMode(t *testing.T) {
	cases := []valueCase{
		{"remote", ModeRemote, false},
		{"test", ModeTest, false},
		{"empty", "", true},
		{"unknown", "local", true},
		{"wrong case", "Remote", true},
	}

	runValueCases(t, cases, func(v string) error {
		return checkMode(&Config{RepositoryMode: v})
	})
}

// TestCheckMaxDepth covers the boundaries of the allowed depth.
func TestCheckMaxDepth(t *testing.T) {
	cases := []struct {
		depth   int
		wantErr bool
	}{
		{minDepth - 1, true},
		{minDepth, false},
		{maxDepthLimit, false},
		{maxDepthLimit + 1, true},
	}

	for _, tc := range cases {
		err := checkMaxDepth(&Config{MaxDepth: tc.depth})

		if (err != nil) != tc.wantErr {
			t.Errorf("depth %d: got error %v, want error: %t",
				tc.depth, err, tc.wantErr)
		}
	}
}

// TestCheckRepositoryURL covers repositories in the remote mode.
func TestCheckRepositoryURL(t *testing.T) {
	cases := []valueCase{
		{"https", "https://registry.npmjs.org", false},
		{"http with path", "http://localhost:8080/repo", false},
		{"empty", "", true},
		{"ftp", "ftp://example.com/repo", true},
		{"no scheme", "example.com/repo", true},
		{"no host", "https://", true},
		{"broken", "https://exa mple.com", true},
	}

	runValueCases(t, cases, func(v string) error {
		cfg := &Config{Repository: v, RepositoryMode: ModeRemote}
		return checkRepository(cfg)
	})
}

// TestCheckRepositoryFile covers repositories in the test mode.
func TestCheckRepositoryFile(t *testing.T) {
	file := writeFile(t, "repo.txt", "A: B\n")
	cases := []valueCase{
		{"existing file", file, false},
		{"missing file", filepath.Join(filepath.Dir(file), "no.txt"), true},
		{"directory", filepath.Dir(file), true},
		{"empty", "", true},
	}

	runValueCases(t, cases, func(v string) error {
		cfg := &Config{Repository: v, RepositoryMode: ModeTest}
		return checkRepository(cfg)
	})
}
