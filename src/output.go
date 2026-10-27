package main

import (
	"fmt"
	"io"

	"depgraph/src/apk"
)

// printDependencies writes the direct dependencies of the package.
func printDependencies(w io.Writer, pkg *apk.Package) {
	deps := pkg.DirectDependencies()

	fmt.Fprintf(w, "Package: %s %s\n", pkg.Name, pkg.Version)

	if len(deps) == 0 {
		fmt.Fprintln(w, "Direct dependencies: none")
		return
	}

	fmt.Fprintf(w, "Direct dependencies (%d):\n", len(deps))

	for _, dep := range deps {
		fmt.Fprintf(w, "  %s\n", dep)
	}
}
