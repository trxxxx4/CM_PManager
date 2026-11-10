package graph

import "depgraph/src/apk"

// BuildReverse finds packages depending on root, directly or transitively.
// Edges point from the searched package to its dependents for DFS traversal.
func BuildReverse(packages []apk.Package, root string, maxDepth int) (*Graph, error) {
	reversed := reversePackages(packages)
	return Build(reversed, root, maxDepth)
}

// reversePackages turns each dependency into a link to the dependent package.
func reversePackages(packages []apk.Package) []apk.Package {
	byName := make(map[string]*apk.Package, len(packages))
	for _, pkg := range packages {
		if _, exists := byName[pkg.Name]; !exists {
			byName[pkg.Name] = &apk.Package{Name: pkg.Name, Version: pkg.Version}
		}
	}
	for _, pkg := range packages {
		for _, dependency := range pkg.DirectDependencies() {
			target := byName[dependency.Name]
			if target == nil {
				continue
			}
			target.Dependencies = append(target.Dependencies, apk.Dependency{Name: pkg.Name})
		}
	}
	result := make([]apk.Package, 0, len(byName))
	for _, pkg := range byName {
		result = append(result, *pkg)
	}
	return result
}
