// NewHandler constructs and returns a useable request handler.

func NewHandler(c HandlerConfig) *Handler {
	if c.RequestDuration == nil {
		c.RequestDuration = metrics.NewNopHistogram()
	}
	if c.Logger == nil {
		c.Logger = log.NewNopLogger()
	}

	return &Handler{
		db:     c.DB,
		dur:    c.RequestDuration,
		logger: c.Logger,

	}
}

// HandlerConfig captures the dependencies used by the Handler.
type HandlerConfig struct {
	// DB is the backing SQL data store. Required.
	DB *sql.DB

	// RequestDuration will receive observations in seconds.
	// Optional; if nil, a no-op histogram will be used.
	RequestDuration *metrics.Histogram

	// Logger is used to log warnings unsuitable for clients.
	// Optional; if nil, a no-op logger will be used.
	Logger *log.Logger
}


// If a component has a few required dependencies and many optional dependencies, the functional options idiom may be a good fit.
// NewHandler constructs and returns a useable request handler.
func NewHandler(db *sql.DB, options ...HandlerOption) *Handler {
	h := &Handler{
		db:     c.DB,
		dur:    metrics.NewNopHistogram(),
		logger: log.NewNopLogger(),
	}
	for _, option := range options {
		option(h)
	}
	return h
}

// HandlerOption sets an option on the Handler.
type HandlerOption func(*Handler)

// WithRequestDuration injects a histogram to receive observations in seconds.
// By default, a no-op histogram will be used.
func WithRequestDuration(dur *metrics.Histogram) HandlerOption {
	return func(h *Handler) { h.dur = dur }
}

// WithLogger injects a logger to log warnings unsuitable for clients.
// By default, a no-op logger will be used.
func WithLogger(logger *log.Logger) HandlerOption {
	return func(h *Handler) { h.logger = logger }
}