package nestdepth

// OverriddenDeep has depth 7 but override raises thresholds.
// if(1) -> for(2) -> switch(3) -> case(4) -> if(5) -> for(6) -> if(7)
// With override warn=8,fail=10: depth 7 is green zone. No diagnostic.
//
//complexity:nestdepth:warn=8,fail=10
func OverriddenDeep() {
	if true {
		for i := 0; i < 10; i++ {
			switch i {
			case 1:
				if true {
					for j := 0; j < 5; j++ {
						if true {
							_ = j
						}
					}
				}
			}
		}
	}
}
