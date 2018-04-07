package client

import (
	"io"

	"github.com/xakep666/wurl/pkg/config"
)

// Client is an interface for websocket operations
type Client interface {
	// ReadTo reads all websocket text/binary message to given writer.
	// Useful for streams reading
	ReadTo(writer io.Writer) error

	// Ping sends ping message to server with given payload
	Ping(payload []byte) error

	// WriteSingleMessage sends message of given type with given payload
	WriteSingleMessage(payload []byte, messageType int) error

	// WriteJSONMessage encodes 'obj' to JSON and sends it to server
	WriteJSONMessage(obj interface{}) error

	io.Closer
}

// Constructor is a type for client constructors
type Constructor func(url string, opts *config.Options) (Client, error)
