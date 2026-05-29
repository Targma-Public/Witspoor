package witspoor

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ConsoleEmitter writes completed Op trees to a writer as human-readable output.
// Intended for local development — pass it to New() to see output immediately.
type ConsoleEmitter struct {
	w io.Writer
}

// NewConsoleEmitter returns an emitter that writes to the given writer.
// Pass nil to write to stdout.
func NewConsoleEmitter(w io.Writer) *ConsoleEmitter {
	if w == nil {
		w = os.Stdout
	}
	return &ConsoleEmitter{w: w}
}

func (e *ConsoleEmitter) Emit(op *Op) error {
	e.writeOp(op, 0)
	fmt.Fprintln(e.w)
	return nil
}

func (e *ConsoleEmitter) writeOp(op *Op, depth int) {
	indent := strings.Repeat("  ", depth)
	health := fmt.Sprintf("[%s]", op.Health.String())
	duration := op.Duration()

	fmt.Fprintf(e.w, "%s%-12s %s  (%s)\n", indent, health, op.Name, duration)

	inner := indent + "  "

	for _, f := range op.facts {
		if len(f.meta) > 0 {
			fmt.Fprintf(e.w, "%sfact        %s = %v  %v\n", inner, f.Name, f.Value, f.Meta())
		} else {
			fmt.Fprintf(e.w, "%sfact        %s = %v\n", inner, f.Name, f.Value)
		}
	}

	for _, r := range op.readings {
		suffix := e.constraintSuffix(r)
		fmt.Fprintf(e.w, "%sreading     %s = %v%s\n", inner, r.Name, r.Value, suffix)
	}

	for _, i := range op.incidents {
		fmt.Fprintf(e.w, "%sincident    %s: %v\n", inner, i.Name, i.Err)
	}

	for _, child := range op.attempts {
		e.writeOp(child, depth+1)
	}

	for _, child := range op.children {
		e.writeOp(child, depth+1)
	}
}

func (e *ConsoleEmitter) constraintSuffix(r *Reading) string {
	if len(r.constraints) == 0 {
		return ""
	}
	parts := make([]string, len(r.constraints))
	for i, c := range r.constraints {
		var level string
		switch c.level {
		case Degraded:
			level = "warn"
		case Failed:
			level = "error"
		default:
			level = c.level.String()
		}
		parts[i] = fmt.Sprintf("%s %s %.4g", level, c.op, c.threshold)
	}
	return "  (" + strings.Join(parts, ", ") + ")"
}
