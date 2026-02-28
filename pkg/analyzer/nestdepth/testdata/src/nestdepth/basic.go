package nestdepth

// Flat has no nesting. Green zone.
func Flat() {
	x := 1
	y := 2
	_ = x + y
}

// Empty has an empty body. Green zone.
func Empty() {}

// DepthOne has nesting depth of 1. Green zone.
func DepthOne() {
	if true {
		_ = 1
	}
}

// DepthFour has nesting depth of 4. Green zone (at boundary).
// if(1) -> for(2) -> switch(3) -> case(4)
func DepthFour() {
	if true {
		for i := 0; i < 10; i++ {
			switch i {
			case 1:
				_ = i
			}
		}
	}
}

// DepthFive has nesting depth of 5. Yellow zone (warning).
// if(1) -> for(2) -> switch(3) -> case(4) -> if(5)
func DepthFive() {
	if true {
		for i := 0; i < 10; i++ {
			switch i {
			case 1:
				if true { // want `function DepthFive has a nesting depth of 5 \(warn: >4, fail: >6\) \[warning\]`
					_ = i
				}
			}
		}
	}
}

// DepthSeven has nesting depth of 7. Red zone (error).
// if(1) -> for(2) -> switch(3) -> case(4) -> if(5) -> for(6) -> if(7)
func DepthSeven() {
	if true {
		for i := 0; i < 10; i++ {
			switch i {
			case 1:
				if true {
					for j := 0; j < 5; j++ {
						if true { // want `function DepthSeven has a nesting depth of 7 \(warn: >4, fail: >6\) \[error\]`
							_ = j
						}
					}
				}
			}
		}
	}
}
