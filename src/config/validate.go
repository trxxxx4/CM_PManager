package config

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	minDepth             = 1
	maxDepthLimit        = 100
	maxPackageNameLength = 214
	maxVersionLength     = 128
)

var (
	packageNamePattern = regexp.MustCompile(`^[A-Za-z0-9@][A-Za-z0-9@._/:+-]*$`)
	versionPattern     = regexp.MustCompile(`^[0-9A-Za-z][0-9A-Za-z.+_~:-]*$`)
	allowedSchemes     = map[string]bool{"http": true, "https": true}
)

// check ties a parameter name to the function validating its value.
type check struct {
	key string
	run func(*Config) error
}

var checks = []check{
	{keyPackageName, checkPackageName},
	{keyRepositoryMode, checkMode},
	{keyRepository, checkRepository},
	{keyVersion, checkVersion},
	{keyMaxDepth, checkMaxDepth},
}

// validate checks the values of the parameters that are defined.
func (c *Config) validate(isDefined func(string) bool) []error {
	var errs []error

	for _, item := range checks {
		if !isDefined(item.key) {
			continue
		}

		if err := item.run(c); err != nil {
			errs = append(errs, &ParamError{Key: item.key, Err: err})
		}
	}

	return errs
}

// checkPackageName validates the package name.
func checkPackageName(c *Config) error {
	name := c.PackageName

	if strings.TrimSpace(name) == "" {
		return errors.New("must not be empty")
	}

	if len(name) > maxPackageNameLength {
		return fmt.Errorf("must be at most %d characters", maxPackageNameLength)
	}

	if !packageNamePattern.MatchString(name) {
		return fmt.Errorf("contains invalid characters: %q", name)
	}

	return nil
}

// checkVersion validates the package version.
func checkVersion(c *Config) error {
	version := c.Version

	if strings.TrimSpace(version) == "" {
		return errors.New("must not be empty")
	}

	if len(version) > maxVersionLength {
		return fmt.Errorf("must be at most %d characters", maxVersionLength)
	}

	if !versionPattern.MatchString(version) {
		return fmt.Errorf("invalid version format: %q", version)
	}

	return nil
}

// checkMode validates the repository mode.
func checkMode(c *Config) error {
	switch c.RepositoryMode {
	case ModeRemote, ModeTest:
		return nil
	}

	return fmt.Errorf(
		"must be %q or %q, got %q",
		ModeRemote,
		ModeTest,
		c.RepositoryMode,
	)
}

// checkMaxDepth validates the maximum depth of the dependency analysis.
func checkMaxDepth(c *Config) error {
	if c.MaxDepth < minDepth {
		return fmt.Errorf("must be at least %d, got %d", minDepth, c.MaxDepth)
	}

	if c.MaxDepth > maxDepthLimit {
		return fmt.Errorf(
			"must be at most %d, got %d",
			maxDepthLimit,
			c.MaxDepth,
		)
	}

	return nil
}

// checkRepository validates the repository according to its mode.
func checkRepository(c *Config) error {
	if strings.TrimSpace(c.Repository) == "" {
		return errors.New("must not be empty")
	}

	switch c.RepositoryMode {
	case ModeRemote:
		return checkURL(c.Repository)
	case ModeTest:
		return checkTestFile(c.Repository)
	}

	return nil
}

// checkURL checks that raw is an absolute http or https URL.
func checkURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid URL %q: %w", raw, err)
	}

	if !allowedSchemes[parsed.Scheme] {
		return fmt.Errorf(
			"URL scheme must be http or https, got %q",
			parsed.Scheme,
		)
	}

	if parsed.Host == "" {
		return fmt.Errorf("URL %q has no host", raw)
	}

	return nil
}

// checkTestFile checks that path is an existing regular file.
func checkTestFile(path string) error {
	info, err := os.Stat(path)

	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("test repository file not found: %s", path)
	}

	if err != nil {
		return fmt.Errorf("cannot access %s: %w", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory, a file is expected", path)
	}

	return nil
}
