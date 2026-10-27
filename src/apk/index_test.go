package apk

import (
	"errors"
	"strings"
	"testing"
)

const (
	sampleCount           = 3
	firstDependencyCount  = 4
	directDependencyCount = 3
	availableVersionCount = 2
)

const sampleIndex = `C:Q1111
P:curl
V:8.5.0-r0
D:ca-certificates so:libc.musl-x86_64.so.1 !curl-old libcurl>=8.5.0
p:cmd:curl=8.5.0-r0

P:curl
V:8.4.0-r0

P:musl
V:1.2.4-r4
`

// TestParseDependency covers the forms of dependency tokens.
func TestParseDependency(t *testing.T) {
	virtual := "so:libc.musl-x86_64.so.1"
	cases := []struct {
		token string
		want  Dependency
	}{
		{"libssl3", Dependency{Name: "libssl3"}},
		{"libssl3>=3.1", Dependency{Name: "libssl3", Constraint: ">=3.1"}},
		{"a<2.0", Dependency{Name: "a", Constraint: "<2.0"}},
		{virtual, Dependency{Name: virtual}},
		{"pc:openssl~3.1", Dependency{Name: "pc:openssl", Constraint: "~3.1"}},
		{"!old", Dependency{Name: "old", Conflict: true}},
	}

	for _, tc := range cases {
		if got := ParseDependency(tc.token); got != tc.want {
			t.Errorf("%q: got %+v, want %+v", tc.token, got, tc.want)
		}
	}
}

// TestDependencyString checks that parsing and printing round-trip.
func TestDependencyString(t *testing.T) {
	for _, token := range []string{"name", "name>=1.0", "!name", "!name=2"} {
		if got := ParseDependency(token).String(); got != token {
			t.Errorf("got %q, want %q", got, token)
		}
	}
}

// TestParseIndex checks the records read from the index text.
func TestParseIndex(t *testing.T) {
	packages, err := ParseIndex(strings.NewReader(sampleIndex))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(packages) != sampleCount {
		t.Fatalf("got %d packages, want %d", len(packages), sampleCount)
	}

	first := packages[0]
	if first.Name != "curl" || first.Version != "8.5.0-r0" {
		t.Errorf("unexpected first package: %+v", first)
	}

	if len(first.Dependencies) != firstDependencyCount {
		t.Errorf("got %d dependencies, want %d",
			len(first.Dependencies), firstDependencyCount)
	}

	if len(packages[1].Dependencies) != 0 {
		t.Errorf("second package must have no dependencies")
	}
}

// TestDirectDependencies checks that conflicts are not dependencies.
func TestDirectDependencies(t *testing.T) {
	packages, err := ParseIndex(strings.NewReader(sampleIndex))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	deps := packages[0].DirectDependencies()
	if len(deps) != directDependencyCount {
		t.Fatalf("got %d direct dependencies, want %d",
			len(deps), directDependencyCount)
	}

	for _, dep := range deps {
		if dep.Conflict {
			t.Errorf("conflict %s must be skipped", dep)
		}
	}
}

// TestFindPackage covers a found package and both kinds of misses.
func TestFindPackage(t *testing.T) {
	packages, err := ParseIndex(strings.NewReader(sampleIndex))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := FindPackage(packages, "curl", "8.4.0-r0")
	if err != nil || found.Version != "8.4.0-r0" {
		t.Errorf("got %v, %v", found, err)
	}

	_, err = FindPackage(packages, "missing", "1.0")
	if !errors.Is(err, ErrPackageNotFound) {
		t.Errorf("got %v, want ErrPackageNotFound", err)
	}

	var versionErr *VersionNotFoundError

	_, err = FindPackage(packages, "curl", "1.0")
	if !errors.As(err, &versionErr) ||
		len(versionErr.Available) != availableVersionCount {
		t.Errorf("got %v, want VersionNotFoundError with %d versions",
			err, availableVersionCount)
	}
}
