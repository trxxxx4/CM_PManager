package apk

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var testGraphLine = regexp.MustCompile(`^([A-Z]):(?:\s*(.*))?$`)

// LoadTestGraph reads a simple test graph: one line per package, e.g. A: B C.
func LoadTestGraph(path, version string) ([]Package, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer file.Close()

	var packages []Package
	seen := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		match := testGraphLine.FindStringSubmatch(line)
		if match == nil {
			return nil, fmt.Errorf("invalid test graph at line %d: expected 'A: B C'", lineNumber)
		}
		name := match[1]
		if seen[name] {
			return nil, fmt.Errorf("duplicate package %q at line %d", name, lineNumber)
		}
		seen[name] = true
		packages = append(packages, Package{
			Name:         name,
			Version:      version,
			Dependencies: parseDependencies(match[2]),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("cannot read test graph: %w", err)
	}
	if len(packages) == 0 {
		return nil, fmt.Errorf("test graph is empty: %s", path)
	}
	return packages, nil
}

// LoadTestRepository detects an APK index or the compact uppercase test graph.
func LoadTestRepository(path, version string) ([]Package, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	text := string(content)
	isArchive := strings.HasPrefix(text, "\x1f\x8b")
	hasPackageField := strings.Contains(text, "\nP:") || strings.HasPrefix(text, "P:")
	isAPKIndex := hasPackageField && strings.Contains(text, "\nV:")
	if isArchive || isAPKIndex {
		return LoadFile(path)
	}
	return LoadTestGraph(path, version)
}
