package apk

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	fieldName         = "P"
	fieldVersion      = "V"
	fieldDepends      = "D"
	initialLineBuffer = 64 * 1024
	maxLineLength     = 4 * 1024 * 1024
)

// ErrPackageNotFound means that the index has no package with the name.
var ErrPackageNotFound = errors.New("package not found")

// Package is a record of the package index.
type Package struct {
	Name         string
	Version      string
	Dependencies []Dependency
}

// VersionNotFoundError means that the package exists in other versions.
type VersionNotFoundError struct {
	Name      string
	Version   string
	Available []string
}

// Error describes the missing version and lists the available ones.
func (e *VersionNotFoundError) Error() string {
	return fmt.Sprintf(
		"version %s of package %s not found, available: %s",
		e.Version,
		e.Name,
		strings.Join(e.Available, ", "),
	)
}

// DirectDependencies returns the dependencies without the conflicts.
func (p *Package) DirectDependencies() []Dependency {
	var deps []Dependency

	for _, dep := range p.Dependencies {
		if !dep.Conflict {
			deps = append(deps, dep)
		}
	}

	return deps
}

// ParseIndex reads the text of an APKINDEX file.
func ParseIndex(r io.Reader) ([]Package, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, initialLineBuffer), maxLineLength)

	var packages []Package
	var current Package

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")

		if line == "" {
			packages = appendPackage(packages, current)
			current = Package{}

			continue
		}

		applyField(&current, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("cannot read index: %w", err)
	}

	return appendPackage(packages, current), nil
}

// FindPackage returns the package with the exact name and version.
func FindPackage(packages []Package, name, version string) (*Package, error) {
	var available []string

	for i := range packages {
		if packages[i].Name != name {
			continue
		}

		if packages[i].Version == version {
			return &packages[i], nil
		}

		available = append(available, packages[i].Version)
	}

	if len(available) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrPackageNotFound, name)
	}

	return nil, &VersionNotFoundError{
		Name:      name,
		Version:   version,
		Available: available,
	}
}

// appendPackage adds the package to the list unless it is empty.
func appendPackage(packages []Package, pkg Package) []Package {
	if pkg.Name == "" {
		return packages
	}

	return append(packages, pkg)
}

// applyField stores the value of one "X:value" line in the package.
func applyField(pkg *Package, line string) {
	key, value, found := strings.Cut(line, ":")
	if !found {
		return
	}

	switch key {
	case fieldName:
		pkg.Name = value
	case fieldVersion:
		pkg.Version = value
	case fieldDepends:
		pkg.Dependencies = parseDependencies(value)
	}
}

// parseDependencies converts a space-separated list to dependencies.
func parseDependencies(value string) []Dependency {
	var deps []Dependency

	for _, token := range strings.Fields(value) {
		deps = append(deps, ParseDependency(token))
	}

	return deps
}
