package fanout

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

// RepeatedCall has fan out of 1 (same function called twice). Green zone.
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
