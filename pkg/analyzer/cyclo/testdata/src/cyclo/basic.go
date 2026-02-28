package cyclo

// Simple has complexity 1 (base only). Green zone.
func Simple() {
	x := 1
	_ = x
}

// SingleIf has complexity 2 (1 + if). Green zone.
func SingleIf(x int) {
	if x > 0 {
		_ = x
	}
}

// ForLoop has complexity 2 (1 + for). Green zone.
func ForLoop() {
	for i := 0; i < 10; i++ {
		_ = i
	}
}

// ElseIfChain has complexity 3 (1 + if + else-if). Green zone.
// else does not count but the second if in "else if" does.
func ElseIfChain(x int) {
	if x > 10 {
		_ = 1
	} else if x > 5 {
		_ = 2
	} else {
		_ = 3
	}
}

// SwitchCases has complexity 4 (1 + case + case + case). Green zone.
// default does not count. switch itself does not count.
func SwitchCases(x int) {
	switch x {
	case 1:
		_ = 1
	case 2:
		_ = 2
	case 3:
		_ = 3
	default:
		_ = 0
	}
}

// ComplexFunc has complexity 10 (yellow zone).
// 1 + if + for + range + case + case + case + if + if + if = 10
func ComplexFunc(x int, items []int) { // want `function ComplexFunc has cyclomatic complexity of 10 \(warn: >9, fail: >14\) \[warning\]`
	if x > 0 {
		for i := 0; i < x; i++ {
			_ = i
		}
		for range items {
			_ = 1
		}
		switch x {
		case 1:
			_ = 1
		case 2:
			_ = 2
		case 3:
			_ = 3
		}
		if x > 5 {
			_ = 5
		}
		if x > 10 {
			_ = 10
		}
		if x > 20 {
			_ = 20
		}
	}
}

// ErrGuardExempt tests that error guard clauses don't count.
// Complexity = 1 (base) + for = 2. The if err != nil is exempt. Green zone.
func ErrGuardExempt() (int, error) {
	for i := 0; i < 10; i++ {
		if err := doSomething(); err != nil {
			return 0, err
		}
	}
	return 1, nil
}

// ErrGuardNonExempt tests that non-guard error checks DO count.
// Complexity = 1 + if = 2. The if has two statements so it's not a guard. Green zone.
func ErrGuardNonExempt() error {
	if err := doSomething(); err != nil {
		log(err)
		return err
	}
	return nil
}

// SelectCases has complexity 3 (1 + comm + comm). Green zone.
// default does not count. select itself does not count.
func SelectCases(ch1, ch2 chan int) {
	select {
	case <-ch1:
		_ = 1
	case <-ch2:
		_ = 2
	default:
		_ = 0
	}
}

// RangeLoop has complexity 2 (1 + range). Green zone.
func RangeLoop(items []int) {
	for _, v := range items {
		_ = v
	}
}

// NestedIf has complexity 4 (1 + if + if + if). Green zone.
// Each if is a separate decision, even when nested.
func NestedIf(x int) {
	if x > 0 {
		if x > 5 {
			if x > 10 {
				_ = x
			}
		}
	}
}

func doSomething() error { return nil }
func log(err error)      {}
