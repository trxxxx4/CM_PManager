package visualize

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// FileNames returns safe, package-specific paths for DOT and PNG artifacts.
func FileNames(root string, reverse bool) (string, string) {
	prefix := "dependency_graph_"
	if reverse {
		prefix = "reverse_dependency_graph_"
	}
	name := safeFileToken(root)
	return prefix + name + ".dot", prefix + name + ".png"
}

// safeFileToken replaces characters that are unsuitable for artifact names.
func safeFileToken(value string) string {
	var output strings.Builder
	for _, char := range value {
		if isFileCharacter(char) {
			output.WriteRune(char)
		} else {
			output.WriteByte('_')
		}
	}
	return output.String()
}

// isFileCharacter reports whether a rune is safe in a portable file name.
func isFileCharacter(char rune) bool {
	return char >= 'a' && char <= 'z' ||
		char >= 'A' && char <= 'Z' ||
		char >= '0' && char <= '9' || char == '-' || char == '_'
}

// GenerateFiles writes a DOT file and renders it into a PNG image with Graphviz.
func GenerateFiles(dotPath, imagePath, source string) error {
	if err := os.WriteFile(dotPath, []byte(source), 0o644); err != nil {
		return fmt.Errorf("cannot write Graphviz file %s: %w", dotPath, err)
	}
	dot, err := exec.LookPath("dot")
	if err != nil {
		return fmt.Errorf("Graphviz is required: executable 'dot' was not found")
	}
	command := exec.Command(dot, "-Tpng", dotPath, "-o", imagePath)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Graphviz rendering failed: %s: %w", strings.TrimSpace(string(output)), err)
	}
	return nil
}

// OpenImage asks the Linux desktop to open the generated image.
func OpenImage(path string) error {
	opener, err := exec.LookPath("xdg-open")
	if err != nil {
		return fmt.Errorf("cannot open image: xdg-open was not found")
	}
	command := exec.Command(opener, path)
	if err := command.Start(); err != nil {
		return fmt.Errorf("cannot open image %s: %w", path, err)
	}
	return command.Process.Release()
}

// HasDesktop reports whether a graphical Linux session is available.
func HasDesktop() bool {
	return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
}
