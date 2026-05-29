package witspoor

// Client is the entrypoint for creating Ops and emitting telemetry.
// Construct one with New and inject it wherever you need instrumentation.
type Client struct {
	emitter Emitter
}

// New creates a Client with the given Emitter.
func New(emitter Emitter) *Client {
	return &Client{emitter: emitter}
}

// Op starts a new root-level Op.
func (c *Client) Op(name string) *Op {
	return newOp(name, c, nil)
}

func (c *Client) emit(op *Op) {
	if c.emitter == nil {
		return
	}
	c.emitter.Emit(op)
}
