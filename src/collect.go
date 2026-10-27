package main

import (
	"net/http"

	"depgraph/src/apk"
	"depgraph/src/config"
)

// findPackage loads the repository index and finds the configured package.
func findPackage(cfg *config.Config) (*apk.Package, error) {
	packages, err := loadPackages(cfg)
	if err != nil {
		return nil, err
	}

	return apk.FindPackage(packages, cfg.PackageName, cfg.Version)
}

// loadPackages reads the index from the file or from the URL.
func loadPackages(cfg *config.Config) ([]apk.Package, error) {
	if cfg.RepositoryMode == config.ModeTest {
		return apk.LoadFile(cfg.Repository)
	}

	client := &http.Client{Timeout: requestTimeout}

	return apk.LoadRemote(cfg.Repository, client)
}
