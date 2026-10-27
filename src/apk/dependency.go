// Package apk reads package information from Alpine Linux repositories.
package apk

import "strings"

const (
	conflictPrefix = "!"
	operatorChars  = "<>=~"
)

// Dependency is one entry of the dependency list of a package.
type Dependency struct {
	Name       string
	Constraint string
	Conflict   bool
}

// ParseDependency splits a token like "libssl3>=3.1" into its parts.
func ParseDependency(token string) Dependency {
	conflict := strings.HasPrefix(token, conflictPrefix)
	token = strings.TrimPrefix(token, conflictPrefix)

	index := strings.IndexAny(token, operatorChars)
	if index < 0 {
		return Dependency{Name: token, Conflict: conflict}
	}

	return Dependency{
		Name:       token[:index],
		Constraint: token[index:],
		Conflict:   conflict,
	}
}

// String returns the dependency in the notation of the index.
func (d Dependency) String() string {
	text := d.Name + d.Constraint

	if d.Conflict {
		return conflictPrefix + text
	}

	return text
}
