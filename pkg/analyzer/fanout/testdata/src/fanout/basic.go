package fanout

import "ext.pkg/dep"

// NoCalls has fan out of 0. Green zone.
func NoCalls() {
	x := 1
	_ = x
}

// BuiltinsOnly has fan out of 0 (builtins excluded). Green zone.
func BuiltinsOnly() {
	s := make([]int, 10)
	s = append(s, 1)
	_ = len(s)
}

// StdlibOnly has fan out of 0 (stdlib excluded). Green zone.
func StdlibOnly() {
	_ = helper()
}

// RepeatedCall has fan out of 0 (same-package calls with no-dot path excluded). Green zone.
func RepeatedCall() {
	helper()
	helper()
}

func helper() int { return 1 }

// TypeConversion has fan out of 0 (type conversions excluded). Green zone.
func TypeConversion() {
	x := 42
	_ = float64(x)
	_ = string(rune(x))
}

// SingleExternal has fan out of 1. Green zone.
func SingleExternal() {
	_ = dep.A()
}

// RepeatedExternal has fan out of 1 (same function called twice). Green zone.
func RepeatedExternal() {
	dep.A()
	dep.A()
}

// HighFanOut has fan out of 7 (7 distinct external calls). Yellow zone (warning).
func HighFanOut() { // want `function HighFanOut has fan out of 7 \(warn: >6, fail: >9\) \[warning\]`
	_ = dep.A()
	_ = dep.B()
	_ = dep.C()
	_ = dep.D()
	_ = dep.E()
	_ = dep.F()
	_ = dep.G()
}

// MethodCall has fan out of 2 (function + method are distinct calls). Green zone.
func MethodCall() {
	s := dep.S{}
	_ = dep.A()
	_ = s.Method()
}
