package witspoor

// Fact is a semantic statement attached to an Op. It does not influence health state.
type Fact struct {
	Name  string
	Value any
	meta  map[string]any
}

func newFact(name string, value any) *Fact {
	return &Fact{Name: name, Value: value}
}

// Set attaches metadata to the Fact. Returns the Fact for chaining.
func (f *Fact) Set(key string, value any) *Fact {
	if f.meta == nil {
		f.meta = make(map[string]any)
	}
	f.meta[key] = value
	return f
}

// Meta returns a copy of the fact's metadata.
func (f *Fact) Meta() map[string]any {
	if f.meta == nil {
		return nil
	}
	out := make(map[string]any, len(f.meta))
	for k, v := range f.meta {
		out[k] = v
	}
	return out
}
