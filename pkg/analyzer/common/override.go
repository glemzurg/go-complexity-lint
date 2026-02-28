package common

import (
	"go/ast"
	"strconv"
	"strings"
)

// ParseOverrides scans the doc comments of a FuncDecl for override directives
// of the form:
//
//	//complexity:metricname:warn=N,fail=M
//
// It returns modified thresholds if overrides are found, or the defaults if not.
func ParseOverrides(funcDecl *ast.FuncDecl, metricName string, defaults Thresholds) Thresholds {
	if funcDecl.Doc == nil {
		return defaults
	}

	prefix := "//complexity:" + metricName + ":"

	for _, comment := range funcDecl.Doc.List {
		text := strings.TrimSpace(comment.Text)
		if !strings.HasPrefix(text, prefix) {
			continue
		}

		overrides := text[len(prefix):]
		result := defaults

		for _, part := range strings.Split(overrides, ",") {
			part = strings.TrimSpace(part)
			kv := strings.SplitN(part, "=", 2)
			if len(kv) != 2 {
				continue
			}
			val, err := strconv.Atoi(strings.TrimSpace(kv[1]))
			if err != nil {
				continue
			}
			switch strings.TrimSpace(kv[0]) {
			case "warn":
				result.WarnAt = val
			case "fail":
				result.FailAt = val
			}
		}

		return result
	}

	return defaults
}
