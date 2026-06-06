package witspoor_test

import (
	"errors"
	"os"

	"github.com/witspoor/witspoor-go"
)

func Example_receiptFlow() {
	w := witspoor.New(witspoor.NewConsoleEmitter(os.Stdout))

	op := w.Op("ReceiptProcess")

	op.Fact("LanguageDetected", "english")

	mistral := op.Op("MistralProcessing")
	calls := 0
	mistral.Attempt(func(attempt *witspoor.Op) error {
		calls++
		if calls < 2 {
			return errors.New("context deadline exceeded")
		}
		attempt.Reading("tokens", 1400).Warn("> 2000").Critical("> 5000")
		attempt.Reading("cost", 0.02).Warn("> 0.05")
		attempt.Fact("Model", "mistral-ocr")
		return nil
	}, witspoor.MaxAttempts(3))
	mistral.End()

	db := op.Op("DatabaseWrite")
	db.Reading("duration_ms", 80).Warn("> 500")
	db.End()

	op.End()
}
