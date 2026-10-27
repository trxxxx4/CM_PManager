// Package config loads, validates and prints the user configuration.
package config

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/BurntSushi/toml"
)

// Repository modes supported by the tool.
const (
	ModeRemote = "remote"
	ModeTest   = "test"
)

const (
	keyPackageName    = "package_name"
	keyRepository     = "repository"
	keyRepositoryMode = "repository_mode"
	keyVersion        = "version"
	keyASCIITree      = "ascii_tree"
	keyMaxDepth       = "max_depth"
)

var requiredKeys = []string{
	keyPackageName,
	keyRepository,
	keyRepositoryMode,
	keyVersion,
	keyASCIITree,
	keyMaxDepth,
}

// Config holds the parameters that the user can configure.
type Config struct {
	PackageName    string `toml:"package_name"`
	Repository     string `toml:"repository"`
	RepositoryMode string `toml:"repository_mode"`
	Version        string `toml:"version"`
	ASCIITree      bool   `toml:"ascii_tree"`
	MaxDepth       int    `toml:"max_depth"`
}

// Load reads the TOML file at path and validates every parameter.
func Load(path string) (*Config, error) {
	var cfg Config

	meta, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, readError(path, err)
	}

	errs := structureErrors(meta)
	errs = append(errs, cfg.validate(definedIn(meta))...)

	if len(errs) > 0 {
		return nil, ValidationErrors(errs)
	}

	return &cfg, nil
}

// Validate checks the values of all parameters of the configuration.
func (c *Config) Validate() error {
	errs := c.validate(func(string) bool { return true })

	if len(errs) > 0 {
		return ValidationErrors(errs)
	}

	return nil
}

// readError turns a file or TOML decoding error into a readable one.
func readError(path string, err error) error {
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("config file not found: %s", path)
	}

	return fmt.Errorf("invalid config file %s: %w", path, err)
}

// definedIn returns a function telling whether a key is set in the file.
func definedIn(meta toml.MetaData) func(string) bool {
	return func(key string) bool {
		return meta.IsDefined(key)
	}
}

// structureErrors reports missing required and unknown parameters.
func structureErrors(meta toml.MetaData) []error {
	var errs []error

	for _, key := range requiredKeys {
		if !meta.IsDefined(key) {
			errs = append(errs, &ParamError{Key: key, Err: errMissingParam})
		}
	}

	for _, key := range meta.Undecoded() {
		errs = append(
			errs,
			&ParamError{Key: key.String(), Err: errUnknownParam},
		)
	}

	return errs
}
