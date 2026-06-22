package common

import (
	"go/build"
	"path/filepath"
	"strings"
	"sync"
)

var stdlibCache sync.Map // map[string]bool

// IsStdlib reports whether pkgPath is a Go standard library package.
// A path is stdlib only when go/build resolves it under GOROOT.
func IsStdlib(pkgPath string) bool {
	switch pkgPath {
	case "", "C":
		return false
	case "unsafe", "builtin":
		return true
	}

	if cached, ok := stdlibCache.Load(pkgPath); ok {
		return cached.(bool)
	}

	bp, err := build.Import(pkgPath, filepath.Join(build.Default.GOROOT, "src"), build.FindOnly)
	if err != nil {
		stdlibCache.Store(pkgPath, false)
		return false
	}

	rel, err := filepath.Rel(build.Default.GOROOT, bp.Dir)
	is := err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
	stdlibCache.Store(pkgPath, is)
	return is
}
