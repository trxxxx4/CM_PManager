package config

import (
	"errors"
	"strings"
)

var (
	errMissingParam = errors.New("required parameter is missing")
	errUnknownParam = errors.New("unknown parameter")
)

// ParamError describes a problem with one configuration parameter.
type ParamError struct {
	Key string
	Err error
}

// Error returns the parameter name followed by the problem.
func (e *ParamError) Error() string {
	return e.Key + ": " + e.Err.Error()
}

// Unwrap returns the underlying problem.
func (e *ParamError) Unwrap() error {
	return e.Err
}

// ValidationErrors lists every problem found in the configuration.
type ValidationErrors []error

// Error returns all problems, one per line.
func (v ValidationErrors) Error() string {
	lines := make([]string, 0, len(v))

	for _, err := range v {
		lines = append(lines, "  - "+err.Error())
	}

	return "invalid configuration:\n" + strings.Join(lines, "\n")
}
