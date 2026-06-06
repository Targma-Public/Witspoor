package witspoor_test

import (
	"errors"
	"testing"

	"github.com/witspoor/witspoor-go"
)

func TestAttempt_Recovered(t *testing.T) {
	w := witspoor.New(nil)
	op := w.Op("MistralProcessing")

	calls := 0
	err := op.Attempt(func(attempt *witspoor.Op) error {
		calls++
		if calls < 2 {
			return errors.New("timeout")
		}
		attempt.Reading("tokens", 1400)
		return nil
	}, witspoor.MaxAttempts(3))

	op.End()

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if op.Health != witspoor.Recovered {
		t.Fatalf("expected Recovered, got %s", op.Health)
	}
}

func TestAttempt_Failed(t *testing.T) {
	w := witspoor.New(nil)
	op := w.Op("MistralProcessing")

	err := op.Attempt(func(attempt *witspoor.Op) error {
		return errors.New("timeout")
	}, witspoor.MaxAttempts(3))

	op.End()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if op.Health != witspoor.Failed {
		t.Fatalf("expected Failed, got %s", op.Health)
	}
}

func TestAttempt_Healthy(t *testing.T) {
	w := witspoor.New(nil)
	op := w.Op("MistralProcessing")

	err := op.Attempt(func(attempt *witspoor.Op) error {
		attempt.Reading("tokens", 1400)
		return nil
	}, witspoor.MaxAttempts(3))

	op.End()

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if op.Health != witspoor.Healthy {
		t.Fatalf("expected Healthy, got %s", op.Health)
	}
}

func TestAttempt_DegradedOnReadingBreach(t *testing.T) {
	w := witspoor.New(nil)
	op := w.Op("MistralProcessing")

	err := op.Attempt(func(attempt *witspoor.Op) error {
		attempt.Reading("tokens", 3000).Warn("> 2000")
		return nil
	})

	op.End()

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if op.Health != witspoor.Degraded {
		t.Fatalf("expected Degraded, got %s", op.Health)
	}
}
