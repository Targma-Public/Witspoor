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
