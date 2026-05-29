package witspoor

// Emitter receives a completed Op tree when a root Op ends.
type Emitter interface {
	Emit(op *Op) error
}

// EmitterFunc is a function that implements Emitter.
type EmitterFunc func(op *Op) error

func (f EmitterFunc) Emit(op *Op) error {
	return f(op)
}

// MultiEmitter fans out to multiple emitters. All emitters are called;
// the first non-nil error is returned.
type MultiEmitter []Emitter

func (m MultiEmitter) Emit(op *Op) error {
	var first error
	for _, e := range m {
		if err := e.Emit(op); err != nil && first == nil {
			first = err
		}
	}
	return first
}

// NoopEmitter silently discards all emissions. Useful in tests.
type NoopEmitter struct{}

func (NoopEmitter) Emit(_ *Op) error { return nil }
