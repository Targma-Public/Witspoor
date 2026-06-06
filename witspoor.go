package witspoor

const source = "witspoor"

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithService tags every emission from this client with a service name.
// Useful for filtering in Axiom or any log aggregator when multiple
// services emit to the same dataset.
func WithService(name string) ClientOption {
	return func(c *Client) { c.service = name }
}

// WithTags adds arbitrary key-value tags stamped on every emission.
// Values must be string-representable.
func WithTags(tags map[string]string) ClientOption {
	return func(c *Client) {
		for k, v := range tags {
			c.tags[k] = v
		}
	}
}

// Client is the entrypoint for creating Ops and emitting telemetry.
// Construct one with New and inject it wherever you need instrumentation.
type Client struct {
	emitter Emitter
	service string
	tags    map[string]string
}

// New creates a Client with the given Emitter and optional options.
func New(emitter Emitter, opts ...ClientOption) *Client {
	c := &Client{
		emitter: emitter,
		tags:    make(map[string]string),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
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
