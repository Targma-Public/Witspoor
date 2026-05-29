package witspoor

// Incident represents an unexpected failure attached to an Op.
// An Incident always influences health state — it is distinct from
// a Reading that breaches a threshold (predicted/capacity problems).
type Incident struct {
	Name string
	Err  error
}

func newIncident(name string, err error) *Incident {
	return &Incident{Name: name, Err: err}
}
