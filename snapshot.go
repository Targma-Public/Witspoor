package witspoor

import "time"

// OpSnapshot is the immutable, serializable form of a completed Op tree.
// This is what emitters receive and what travels over the wire.
type OpSnapshot struct {
	Source     string            `json:"source"`
	Service    string            `json:"service,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
	Name       string            `json:"name"`
	StartedAt  time.Time         `json:"started_at"`
	EndedAt    time.Time         `json:"ended_at"`
	DurationMs float64           `json:"duration_ms"`
	Health     string            `json:"health"`
	Facts      []FactSnapshot    `json:"facts,omitempty"`
	Readings   []ReadingSnapshot `json:"readings,omitempty"`
	Incidents  []IncidentSnapshot `json:"incidents,omitempty"`
	Attempts   []OpSnapshot      `json:"attempts,omitempty"`
	Children   []OpSnapshot      `json:"children,omitempty"`
}

type FactSnapshot struct {
	Name  string         `json:"name"`
	Value any            `json:"value"`
	Meta  map[string]any `json:"meta,omitempty"`
}

type ReadingSnapshot struct {
	Name        string               `json:"name"`
	Value       float64              `json:"value"`
	Constraints []ConstraintSnapshot `json:"constraints,omitempty"`
}

type ConstraintSnapshot struct {
	Level     string  `json:"level"`
	Op        string  `json:"op"`
	Threshold float64 `json:"threshold"`
}

type IncidentSnapshot struct {
	Name  string `json:"name"`
	Error string `json:"error"`
}

// Snapshot converts a completed Op tree into an OpSnapshot.
func (o *Op) Snapshot() OpSnapshot {
	snap := OpSnapshot{
		Source:     source,
		Name:       o.Name,
		StartedAt:  o.StartedAt,
		EndedAt:    o.EndedAt,
		DurationMs: float64(o.Duration().Nanoseconds()) / 1e6,
		Health:     o.Health.String(),
	}

	if o.client != nil {
		snap.Service = o.client.service
		if len(o.client.tags) > 0 {
			snap.Tags = make(map[string]string, len(o.client.tags))
			for k, v := range o.client.tags {
				snap.Tags[k] = v
			}
		}
	}

	for _, f := range o.facts {
		snap.Facts = append(snap.Facts, FactSnapshot{
			Name:  f.Name,
			Value: f.Value,
			Meta:  f.Meta(),
		})
	}

	for _, r := range o.readings {
		rs := ReadingSnapshot{Name: r.Name, Value: r.Value}
		for _, c := range r.constraints {
			level := "degraded"
			if c.level == Failed {
				level = "failed"
			}
			rs.Constraints = append(rs.Constraints, ConstraintSnapshot{
				Level:     level,
				Op:        c.op,
				Threshold: c.threshold,
			})
		}
		snap.Readings = append(snap.Readings, rs)
	}

	for _, i := range o.incidents {
		snap.Incidents = append(snap.Incidents, IncidentSnapshot{
			Name:  i.Name,
			Error: i.Err.Error(),
		})
	}

	for _, a := range o.attempts {
		snap.Attempts = append(snap.Attempts, a.Snapshot())
	}

	for _, c := range o.children {
		snap.Children = append(snap.Children, c.Snapshot())
	}

	return snap
}
