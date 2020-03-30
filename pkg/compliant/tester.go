package compliant

func NewTester() *tester {
	return &tester{}
}

type ListenerFunc func(result ManifestResult) error

type tester struct {
	listeners []ListenerFunc
}

func (t *tester) AddListener(listener ListenerFunc) {
	t.listeners = append(t.listeners, listener)
}

func (t *tester) Compliant(manifest Manifest) (bool, error) {
	result := manifest.Run()

	for _, listener := range t.listeners {
		err := listener(result)
		if err != nil {
			return false, err
		}
	}

	return !result.Fail(), nil
}
