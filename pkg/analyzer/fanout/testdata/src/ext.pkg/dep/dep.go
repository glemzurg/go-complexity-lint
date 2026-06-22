package dep

func A() int { return 1 }
func B() int { return 2 }
func C() int { return 3 }
func D() int { return 4 }
func E() int { return 5 }
func F() int { return 6 }
func G() int { return 7 }
func H() int { return 8 }

func Fail() error { return nil }

func Wrap(err error) error { return err }

type S struct{}

func (s S) Method() int { return 9 }
