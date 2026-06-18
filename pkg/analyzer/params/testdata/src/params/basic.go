package params

import "context"

// NoParams has 0 params. Green zone.
func NoParams() {}

// FourParams has 4 params. Green zone (at boundary).
func FourParams(a, b, c, d int) {
	_ = a + b + c + d
}

// FiveParams has 5 params. Yellow zone (warning).
func FiveParams(a, b, c, d, e int) { // want `function FiveParams has 5 parameters \(warn: >=5, fail: >=7\) \[warning\] \(reduce by grouping coherently related subsets of parameters into structs — do not simply wrap all params into a single struct\)`
	_ = a + b + c + d + e
}

// SevenParams has 7 params. Red zone (error).
func SevenParams(a, b int, c string, d, e, f float64, g bool) { // want `function SevenParams has 7 parameters \(warn: >=5, fail: >=7\) \[error\] \(reduce by grouping coherently related subsets of parameters into structs — do not simply wrap all params into a single struct\)`
	_, _, _, _, _, _, _ = a, b, c, d, e, f, g
}

// GroupedParams has 5 params despite only 2 field entries. Yellow zone.
func GroupedParams(a, b, c int, d, e string) { // want `function GroupedParams has 5 parameters \(warn: >=5, fail: >=7\) \[warning\] \(reduce by grouping coherently related subsets of parameters into structs — do not simply wrap all params into a single struct\)`
	_, _, _, _, _ = a, b, c, d, e
}

// MethodReceiver tests that the receiver does not count as a parameter.
type MyStruct struct{}

func (m *MyStruct) MethodWithFourParams(a, b, c, d int) {
	_ = a + b + c + d
}

// Variadic tests that variadic params count as 1.
func Variadic(a int, b ...string) {
	_, _ = a, b
}

// UnnamedParams tests that unnamed parameters each count as 1.
// 5 unnamed params. Yellow zone (warning).
func UnnamedParams(int, string, bool, float64, error) { // want `function UnnamedParams has 5 parameters \(warn: >=5, fail: >=7\) \[warning\] \(reduce by grouping coherently related subsets of parameters into structs — do not simply wrap all params into a single struct\)`
}

// MixedNamedUnnamed has 5 params (3 named + 2 unnamed fields = 5). Yellow zone.
func MixedNamedUnnamed(a, b int, _ string, _ bool, c float64) { // want `function MixedNamedUnnamed has 5 parameters \(warn: >=5, fail: >=7\) \[warning\] \(reduce by grouping coherently related subsets of parameters into structs — do not simply wrap all params into a single struct\)`
	_, _, _ = a, b, c
}

// WithCtx has 7 listed params but ctx context.Context is exempt, so 6 count. Yellow zone.
func WithCtx(ctx context.Context, a, b, c, d, e, f int) { // want `function WithCtx has 6 parameters \(warn: >=5, fail: >=7\) \[warning\] \(reduce by grouping coherently related subsets of parameters into structs — do not simply wrap all params into a single struct\)`
	_, _, _, _, _, _ = a, b, c, d, e, f
}

// WithCtxGreen has ctx plus 4 others = 4 counted. Green zone.
func WithCtxGreen(ctx context.Context, a, b, c, d int) {
	_, _, _, _ = a, b, c, d
}

// WithCtxRed has ctx plus 7 others = 7 counted. Red zone.
func WithCtxRed(ctx context.Context, a, b, c, d, e, f, g int) { // want `function WithCtxRed has 7 parameters \(warn: >=5, fail: >=7\) \[error\] \(reduce by grouping coherently related subsets of parameters into structs — do not simply wrap all params into a single struct\)`
	_, _, _, _, _, _, _ = a, b, c, d, e, f, g
}

// WrongCtxName has c context.Context, not ctx; all 7 count. Red zone.
func WrongCtxName(c context.Context, a, b, d, e, f, g int) { // want `function WrongCtxName has 7 parameters \(warn: >=5, fail: >=7\) \[error\] \(reduce by grouping coherently related subsets of parameters into structs — do not simply wrap all params into a single struct\)`
	_, _, _, _, _, _ = a, b, d, e, f, g
}

// WrongCtxType has ctx with non-Context type; all 7 count. Red zone.
func WrongCtxType(ctx context.CancelFunc, a, b, c, d, e, f int) { // want `function WrongCtxType has 7 parameters \(warn: >=5, fail: >=7\) \[error\] \(reduce by grouping coherently related subsets of parameters into structs — do not simply wrap all params into a single struct\)`
	_, _, _, _, _, _ = a, b, c, d, e, f
}
