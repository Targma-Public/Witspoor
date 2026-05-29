package witspoor

import (
	"fmt"
	"strconv"
	"strings"
)

// Reading is a numeric measurement attached to an Op.
// Without constraints it is purely informational.
// Breached constraints influence the Op's health state.
type Reading struct {
	Name        string
	Value       float64
	constraints []constraint
}

type constraint struct {
	level HealthState
	op    string  // ">" or "<"
	threshold float64
}

func newReading(name string, value float64) *Reading {
	return &Reading{Name: name, Value: value}
}

// Warn adds a warning-level constraint. Expr format: "> 2000" or "< 100".
func (r *Reading) Warn(expr string) *Reading {
	return r.addConstraint(Degraded, expr)
}

// Critical adds a critical-level constraint. Expr format: "> 5000" or "< 10".
func (r *Reading) Critical(expr string) *Reading {
	return r.addConstraint(Failed, expr)
}

func (r *Reading) addConstraint(level HealthState, expr string) *Reading {
	c, err := parseConstraint(level, expr)
	if err != nil {
		// invalid constraint expression — silently skip rather than panic
		return r
	}
	r.constraints = append(r.constraints, c)
	return r
}

// evaluate returns the worst health state triggered by this reading's value.
func (r *Reading) evaluate() HealthState {
	state := Healthy
	for _, c := range r.constraints {
		if c.breached(r.Value) {
			state = worse(state, c.level)
		}
	}
	return state
}

func (c constraint) breached(value float64) bool {
	switch c.op {
	case ">":
		return value > c.threshold
	case "<":
		return value < c.threshold
	}
	return false
}

func parseConstraint(level HealthState, expr string) (constraint, error) {
	parts := strings.Fields(expr)
	if len(parts) != 2 {
		return constraint{}, fmt.Errorf("witspoor: invalid constraint %q", expr)
	}
	op := parts[0]
	if op != ">" && op != "<" {
		return constraint{}, fmt.Errorf("witspoor: unsupported operator %q in constraint %q", op, expr)
	}
	threshold, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return constraint{}, fmt.Errorf("witspoor: invalid threshold in constraint %q: %w", expr, err)
	}
	return constraint{level: level, op: op, threshold: threshold}, nil
}
