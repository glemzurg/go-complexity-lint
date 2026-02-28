package nestdepth

import "fmt"

// MyStruct is used to test method receiver reporting.
type MyStruct struct{}

// DeepMethod tests method receiver name formatting.
// if(1) -> for(2) -> switch(3) -> case(4) -> select(5) -> comm(6)
// Depth 6 = yellow zone (warning).
func (m *MyStruct) DeepMethod() {
	if true {
		for i := 0; i < 10; i++ {
			switch i {
			case 0:
				select {
				case <-make(chan int):
					_ = i // want `function \*MyStruct\.DeepMethod has a nesting depth of 6 \(warn: >4, fail: >6\) \[warning\]`
				}
			}
		}
	}
}

// ClosureNesting tests that func literals add a nesting level.
// if(1) -> funclit(2) -> for(3) -> switch(4) -> case(5) -> if(6)
// Depth 6 = yellow zone (warning).
func ClosureNesting() {
	if true {
		fn := func() {
			for i := 0; i < 10; i++ {
				switch i {
				case 1:
					if true { // want `function ClosureNesting has a nesting depth of 6 \(warn: >4, fail: >6\) \[warning\]`
						_ = i
					}
				}
			}
		}
		fn()
	}
}

// ElseIfChain tests that else-if does not double-increment.
// Max depth is 1 (each branch is depth 1). Green zone.
func ElseIfChain() {
	if true {
		_ = 1
	} else if false {
		_ = 2
	} else {
		_ = 3
	}
}

// RangeLoop tests range loops.
// range(1) -> if(2) -> switch(3) -> case(4) -> range(5) -> if(6)
// Depth 6 = yellow zone (warning).
func RangeLoop() {
	for _, v := range []int{1, 2, 3} {
		if v > 0 {
			switch v {
			case 1:
				for _, w := range []int{4, 5} {
					if w > 0 { // want `function RangeLoop has a nesting depth of 6 \(warn: >4, fail: >6\) \[warning\]`
						_ = w
					}
				}
			}
		}
	}
}

// TypeSwitchNesting tests type switch. Depth is 3. Green zone.
func TypeSwitchNesting() {
	var x any = 42
	switch x.(type) {
	case int:
		if true {
			_ = 1
		}
	}
}

// DeferClosure tests defer with a func literal. Depth is 2. Green zone.
func DeferClosure() {
	defer func() {
		if true {
			_ = 1
		}
	}()
}

// GoRoutineClosure tests go with a deeply nested func literal.
// funclit(1) -> for(2) -> switch(3) -> case(4) -> select(5) -> default(6) -> if(7)
// Depth 7 = red zone (error).
func GoRoutineClosure() {
	go func() {
		for i := 0; i < 10; i++ {
			switch i {
			case 1:
				select {
				default:
					if true { // want `function GoRoutineClosure has a nesting depth of 7 \(warn: >4, fail: >6\) \[error\]`
						fmt.Println(i)
					}
				}
			}
		}
	}()
}

// ErrGuardExempt tests that error guard clauses don't count as nesting.
// The if err != nil is exempt, so depth is only 1 (the for loop). Green zone.
func ErrGuardExempt() (int, error) {
	for i := 0; i < 10; i++ {
		if err := doSomething(); err != nil {
			return 0, err
		}
	}
	return 1, nil
}

// LabeledLoop tests that labeled statements don't add nesting depth.
// The label itself is transparent; depth comes from the for and nested ifs.
// for(1) -> for(2) -> if(3) -> if(4) -> if(5)
// Depth 5 = yellow zone (warning).
func LabeledLoop() {
outer:
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i > 0 {
				if j > 0 {
					if i+j > 15 { // want `function LabeledLoop has a nesting depth of 5 \(warn: >4, fail: >6\) \[warning\]`
						break outer
					}
				}
			}
		}
	}
}

// BareBlock tests that bare block statements don't add nesting depth.
// The block is transparent; depth comes from the constructs inside it.
// if(1) -> for(2) -> switch(3) -> case(4) -> if(5)
// Depth 5 = yellow zone (warning).
func BareBlock() {
	{
		if true {
			for i := 0; i < 10; i++ {
				switch i {
				case 1:
					if true { // want `function BareBlock has a nesting depth of 5 \(warn: >4, fail: >6\) \[warning\]`
						_ = i
					}
				}
			}
		}
	}
}

func doSomething() error { return nil }
