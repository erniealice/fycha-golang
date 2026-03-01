package fycha

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// CopyStyles copies fycha's CSS assets to the target directory.
// Uses runtime.Caller(0) via packageDir() to discover fycha's package
// directory, same approach as centymo and entydad.
//
// Files are copied to {targetDir}/fycha/ to keep them namespaced.
func CopyStyles(targetDir string) error {
	dir := packageDir()
	if dir == "" {
		return fmt.Errorf("could not determine fycha package directory")
	}

	srcDir := filepath.Join(dir, "assets", "css")
	dstDir := filepath.Join(targetDir, "fycha")

	copied, err := copyDirFiles(srcDir, dstDir, "*.css")
	if err != nil {
		return fmt.Errorf("failed to copy fycha styles: %w", err)
	}

	if copied == 0 {
		log.Printf("fycha: no CSS files found in %s", srcDir)
		return nil
	}

	log.Printf("Copied %d fycha styles to: %s", copied, dstDir)
	return nil
}

// CopyStaticAssets copies fycha's JavaScript assets to the target directory.
// Files are copied to {targetDir}/fycha/ to keep them namespaced.
func CopyStaticAssets(targetDir string) error {
	dir := packageDir()
	if dir == "" {
		return fmt.Errorf("could not determine fycha package directory")
	}

	srcDir := filepath.Join(dir, "assets", "js")
	dstDir := filepath.Join(targetDir, "fycha")

	copied, err := copyDirFiles(srcDir, dstDir, "*.js")
	if err != nil {
		return fmt.Errorf("failed to copy fycha assets: %w", err)
	}

	if copied == 0 {
		log.Printf("fycha: no JS files found in %s", srcDir)
		return nil
	}

	log.Printf("Copied %d fycha assets to: %s", copied, dstDir)
	return nil
}

// copyDirFiles copies all files matching a glob pattern from srcDir to dstDir.
func copyDirFiles(srcDir, dstDir, pattern string) (int, error) {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create target directory: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(srcDir, pattern))
	if err != nil {
		return 0, fmt.Errorf("failed to list source files: %w", err)
	}

	var copied int
	for _, srcFile := range files {
		data, err := os.ReadFile(srcFile)
		if err != nil {
			log.Printf("Warning: Failed to read %s: %v", srcFile, err)
			continue
		}

		dstFile := filepath.Join(dstDir, filepath.Base(srcFile))
		if err := os.WriteFile(dstFile, data, 0644); err != nil {
			return copied, fmt.Errorf("failed to write %s: %w", dstFile, err)
		}
		copied++
	}

	return copied, nil
}
