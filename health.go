package witspoor

type HealthState int

const (
	Healthy   HealthState = iota
	Recovered             // succeeded after at least one error
	Degraded              // reading constraint breached, no unrecovered errors
	Failed                // unrecovered error
)

func (h HealthState) String() string {
	switch h {
	case Healthy:
		return "healthy"
	case Recovered:
		return "recovered"
	case Degraded:
		return "degraded"
	case Failed:
		return "failed"
	default:
		return "unknown"
	}
}

// worse returns the more severe of two health states.
func worse(a, b HealthState) HealthState {
	if b > a {
		return b
	}
	return a
}
