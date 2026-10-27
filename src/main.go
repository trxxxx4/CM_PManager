// Command depgraph visualizes the dependency graph of a package.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"depgraph/src/config"
)

const (
	defaultConfigPath = "config.toml"
	exitOK            = 0
	exitConfigError   = 1
	exitUsageError    = 2
)

// main runs the tool and exits with its status code.
func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run executes the tool and returns the process exit code.
func run(args []string, stdout, stderr io.Writer) int {
	path, err := parseArgs(args, stderr)
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

	cfg.Print(stdout)

	return exitOK
}

// parseArgs reads the command line and returns the config file path.
func parseArgs(args []string, output io.Writer) (string, error) {
	flags := flag.NewFlagSet("depgraph", flag.ContinueOnError)
	flags.SetOutput(output)
	path := flags.String(
		"config",
		defaultConfigPath,
		"path to the TOML configuration file",
	)

	if err := flags.Parse(args); err != nil {
		return "", err
	}

	if flags.NArg() > 0 {
		return "", fmt.Errorf("unexpected arguments: %v", flags.Args())
	}

	return *path, nil
}
