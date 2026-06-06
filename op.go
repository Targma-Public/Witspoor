package witspoor

import "time"

// Op is a bounded operation with start/end/duration.
// It owns its health state, which is inferred at End() from its
// readings, incidents, and children — never set manually.
type Op struct {
	Name      string
	StartedAt time.Time
	EndedAt   time.Time
	Health    HealthState

	facts     []*Fact
	readings  []*Reading
	incidents []*Incident
	children  []*Op
	attempts  []*Op // ordered retry sequence — populated by Attempt()

	client *Client
	parent *Op
}

func newOp(name string, client *Client, parent *Op) *Op {
	return &Op{
		Name:      name,
		StartedAt: time.Now(),
		client:    client,
		parent:    parent,
	}
}

// Op starts a child Op nested under this one.
func (o *Op) Op(name string) *Op {
	child := newOp(name, o.client, o)
	o.children = append(o.children, child)
	return child
}

// Fact attaches a semantic statement to this Op. Does not influence health.
func (o *Op) Fact(name string, value any) *Fact {
	f := newFact(name, value)
	o.facts = append(o.facts, f)
	return f
}

// Reading attaches a numeric measurement to this Op.
func (o *Op) Reading(name string, value float64) *Reading {
	r := newReading(name, value)
	o.readings = append(o.readings, r)
	return r
}

// Incident records an unexpected failure on this Op.
func (o *Op) Incident(name string, err error) *Incident {
	i := newIncident(name, err)
	o.incidents = append(o.incidents, i)
	return i
}

// End closes the Op, infers its health state, and emits the tree
// if this is a root Op (no parent).
func (o *Op) End() {
	o.EndedAt = time.Now()
	o.Health = o.inferHealth()

	if o.parent == nil && o.client != nil {
		o.client.emit(o)
	}
}

// Duration returns how long the Op ran. Zero if not yet ended.
func (o *Op) Duration() time.Duration {
	if o.EndedAt.IsZero() {
		return 0
	}
	return o.EndedAt.Sub(o.StartedAt)
}

// inferHealth derives the health state from readings, incidents, children, and attempts.
// Evaluation order:
//  1. Own incidents → Failed
//  2. Attempt sequence: last succeeded but earlier failed → Recovered; last failed → Failed
//  3. Worst regular child health propagates up
//  4. Reading constraint breaches → Degraded or Failed
func (o *Op) inferHealth() HealthState {
	if len(o.incidents) > 0 {
		return Failed
	}

	if len(o.attempts) > 0 {
		last := o.attempts[len(o.attempts)-1]
		if last.Health == Failed {
			return Failed
		}
		for _, a := range o.attempts[:len(o.attempts)-1] {
			if a.Health == Failed {
				return Recovered
			}
		}
		if last.Health == Degraded {
			return Degraded
		}
	}

	childHealth := Healthy
	for _, child := range o.children {
		childHealth = worse(childHealth, child.Health)
	}

	if childHealth == Failed || childHealth == Recovered {
		return childHealth
	}

	readingHealth := o.evaluateReadings()
	if readingHealth != Healthy {
		return readingHealth
	}

	return childHealth
}

func (o *Op) evaluateReadings() HealthState {
	state := Healthy
	for _, r := range o.readings {
		state = worse(state, r.evaluate())
	}
	return state
}
