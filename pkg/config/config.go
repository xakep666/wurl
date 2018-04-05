package config

import (
	"net/http"
	"time"
)

// Options describes global wurl options
type Options struct {
	// AdditionalHeaders is an additional headers which will be included to request to server
	AdditionalHeaders http.Header

	// AllowInsecureSSL allows to establish insecure SSL connections
	AllowInsecureSSL bool

	// PingPeriod is a period of "ping" messages sending
	PingPeriod time.Duration

	// RespondPings controls whether client will be response with "pong" messages on accepted "ping" messages
	RespondPings bool

	// TraceTo determines where to write debug messages.
	// If TraceTo is empty string debug messages will be disabled.
	// If TraceTo is "-" debug messages will be written to stdout.
	// In other cases debug messages will be written to file.
	TraceTo string
}
