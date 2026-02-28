package params

// NoParams has 0 params. Green zone.
func NoParams() {}

// FourParams has 4 params. Green zone (at boundary).
func FourParams(a, b, c, d int) {
	_ = a + b + c + d
}

// FiveParams has 5 params. Yellow zone (warning).
func FiveParams(a, b, c, d, e int) { // want `function FiveParams has 5 parameters \(warn: >4, fail: >6\) \[warning\]`
	_ = a + b + c + d + e
}

// SevenParams has 7 params. Red zone (error).
func SevenParams(a, b int, c string, d, e, f float64, g bool) { // want `function SevenParams has 7 parameters \(warn: >4, fail: >6\) \[error\]`
	_, _, _, _, _, _, _ = a, b, c, d, e, f, g
}

// GroupedParams has 5 params despite only 2 field entries. Yellow zone.
func GroupedParams(a, b, c int, d, e string) { // want `function GroupedParams has 5 parameters \(warn: >4, fail: >6\) \[warning\]`
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
func UnnamedParams(int, string, bool, float64, error) { // want `function UnnamedParams has 5 parameters \(warn: >4, fail: >6\) \[warning\]`
}

// MixedNamedUnnamed has 5 params (3 named + 2 unnamed fields = 5). Yellow zone.
func MixedNamedUnnamed(a, b int, _ string, _ bool, c float64) { // want `function MixedNamedUnnamed has 5 parameters \(warn: >4, fail: >6\) \[warning\]`
	_, _, _ = a, b, c
}
