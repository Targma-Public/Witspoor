package witspoor

import (
	"fmt"
	"time"
)

// AttemptOption configures the behavior of Op.Attempt.
type AttemptOption func(*attemptConfig)

type attemptConfig struct {
	maxAttempts int
	backoff     time.Duration
}

// MaxAttempts sets the maximum number of attempts. Defaults to 1.
func MaxAttempts(n int) AttemptOption {
	return func(c *attemptConfig) { c.maxAttempts = n }
}

// Backoff sets a fixed wait duration between attempts.
func Backoff(d time.Duration) AttemptOption {
	return func(c *attemptConfig) { c.backoff = d }
}

// Attempt runs fn up to the configured number of times, creating a child Op
// for each attempt. If fn returns an error, the attempt is recorded as Failed
// and the next attempt begins after any configured backoff.
//
// The parent Op's health is inferred from the attempt sequence:
//   - Last attempt succeeded, earlier ones failed → Recovered
//   - Last attempt failed → Failed
//
// The child Op passed to fn is yours to use — attach Readings, Facts, etc.
// Do not call End() on it; Attempt manages that.
func (o *Op) Attempt(fn func(attempt *Op) error, opts ...AttemptOption) error {
	cfg := &attemptConfig{maxAttempts: 1}
	for _, opt := range opts {
		opt(cfg)
	}

	var lastErr error

	for i := 0; i < cfg.maxAttempts; i++ {
		child := newOp(fmt.Sprintf("Attempt %d", i+1), o.client, o)

		lastErr = fn(child)
		if lastErr != nil {
			child.Incident("error", lastErr)
		}

		child.End()
		o.attempts = append(o.attempts, child)

		if lastErr == nil {
			break
		}

		if cfg.backoff > 0 && i < cfg.maxAttempts-1 {
			time.Sleep(cfg.backoff)
		}
	}

	return lastErr
}
