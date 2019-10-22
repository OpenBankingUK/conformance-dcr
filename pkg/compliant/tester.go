package compliant

type Tester interface {
	Compliant(manifest Manifest) (bool, error)
}

func NewTester(expression string, debug bool) Tester {
	tester := NewColourTester(debug)
	if expression != "" {
		return NewFilteredTester(expression, tester)
	}
	return tester
}
