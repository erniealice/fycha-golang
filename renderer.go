package fycha

import (
	"path/filepath"
	"runtime"
)

// TemplatePatterns returns glob patterns for fycha's templates.
// Uses runtime.Caller(0) to discover fycha's package directory,
// same approach as pyeza-golang, centymo, and entydad.
// Consumer apps merge these patterns with pyeza + app patterns when
// initializing the renderer.
func TemplatePatterns() []string {
	dir := packageDir()
	return []string{
		filepath.Join(dir, "templates", "reports", "*.html"),
	}
}

// packageDir returns the absolute directory of this source file.
func packageDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(filename)
}
