package config

import (
	"io"
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
	TraceTo io.WriteCloser

	// ShowHandshakeResponse allows to include handshake response (headers+body) to output
	ShowHandshakeResponse bool

	// Output is an output location for messages
	Output io.WriteCloser

	// ForceBinaryToStdout allows to ignore warning about binary output to stdout
	ForceBinaryToStdout bool

	// MessageAfterConnect allows to send message to server after successful connection.
	MessageAfterConnect io.ReadCloser
}
