package common

import (
	"path/filepath"
	"strings"
)

// ExcludePatterns is a comma-separated list of filename glob patterns.
// Files whose base name matches any pattern are skipped during analysis.
// All analyzers register flags pointing to this variable so that setting
// the flag on any one analyzer applies globally.
var ExcludePatterns string

// IsExcluded reports whether filename's base name matches any of the
// configured exclusion glob patterns.
func IsExcluded(filename string) bool {
	if ExcludePatterns == "" {
		return false
	}
	base := filepath.Base(filename)
	for _, p := range strings.Split(ExcludePatterns, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if matched, _ := filepath.Match(p, base); matched {
			return true
		}
	}
	return false
}
