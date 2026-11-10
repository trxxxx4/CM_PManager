// Command depgraph visualizes the dependency graph of a package.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"depgraph/src/apk"
	"depgraph/src/config"
	"depgraph/src/graph"
)

const (
	defaultConfigPath = "config.toml"
	requestTimeout    = 60 * time.Second
	exitOK            = 0
	exitConfigError   = 1
	exitUsageError    = 2
	exitDataError     = 3
)

// main runs the tool and exits with its status code.
func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run executes the tool and returns the process exit code.
func run(args []string, stdout, stderr io.Writer) int {
	path, reverse, err := parseArgs(args, stderr)
	if errors.Is(err, flag.ErrHelp) {
		return exitOK
	}

	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return exitUsageError
	}

	cfg, err := config.Load(path)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return exitConfigError
	}

	return execute(cfg, reverse, stdout, stderr)
}

// execute loads package data, builds the graph and prints the result.
func execute(cfg *config.Config, reverse bool, stdout, stderr io.Writer) int {
	packages, err := loadPackages(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return exitDataError
	}
	pkg, err := apk.FindPackage(packages, cfg.PackageName, cfg.Version)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return exitDataError
	}
	result, err := buildGraph(packages, pkg.Name, cfg.MaxDepth, reverse)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return exitDataError
	}
	printGraph(stdout, result, reverse)
	return exitOK
}

// parseArgs reads the command line and returns the config file path.
func parseArgs(args []string, output io.Writer) (string, bool, error) {
	flags := flag.NewFlagSet("depgraph", flag.ContinueOnError)
	flags.SetOutput(output)
	path := flags.String(
		"config",
		defaultConfigPath,
		"path to the TOML configuration file",
	)

	reverse := flags.Bool("reverse", false, "show reverse dependencies")
	if err := flags.Parse(args); err != nil {
		return "", false, err
	}

	if flags.NArg() > 0 {
		return "", false, fmt.Errorf("unexpected arguments: %v", flags.Args())
	}

	return *path, *reverse, nil
}

// buildGraph selects the traversal direction without changing stage 3 DFS.
func buildGraph(packages []apk.Package, name string, depth int, reverse bool) (*graph.Graph, error) {
	if reverse {
		return graph.BuildReverse(packages, name, depth)
	}
	return graph.Build(packages, name, depth)
}
