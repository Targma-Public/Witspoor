package witspoor

import (
	"encoding/json"
	"io"
	"os"
)

// JSONEmitter writes completed Op trees as newline-delimited JSON.
// Each emission is one JSON object followed by a newline — suitable for
// log aggregators, file ingestion, and the witspoor backend.
type JSONEmitter struct {
	w       io.Writer
	encoder *json.Encoder
}

// NewJSONEmitter returns an emitter that writes to the given writer.
// Pass nil to write to stdout.
func NewJSONEmitter(w io.Writer) *JSONEmitter {
	if w == nil {
		w = os.Stdout
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &JSONEmitter{w: w, encoder: enc}
}

// NewPrettyJSONEmitter returns a JSONEmitter with indented output.
// Useful for debugging — not recommended for production ingestion.
func NewPrettyJSONEmitter(w io.Writer) *JSONEmitter {
	if w == nil {
		w = os.Stdout
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return &JSONEmitter{w: w, encoder: enc}
}

func (e *JSONEmitter) Emit(op *Op) error {
	return e.encoder.Encode(op.Snapshot())
}
