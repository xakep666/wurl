package client

import (
	"io"

	"github.com/xakep666/wurl/pkg/config"
)

type MessageType int

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage MessageType = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage MessageType = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage MessageType = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage MessageType = 9

	// PongMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage MessageType = 10
)

// Client is an interface for websocket operations
type Client interface {
	// ReadTo reads all websocket text/binary message to given writer.
	// Useful for streams reading
	ReadTo(writer io.Writer) error

	// Ping sends ping message to server with given payload
	Ping(payload []byte) error

	// WriteSingleMessage sends message of given type with given payload
	WriteSingleMessage(payload []byte, messageType MessageType) error

	// WriteJSONMessage encodes 'obj' to JSON and sends it to server
	WriteJSONMessage(obj interface{}) error

	io.Closer
}

// Constructor is a type for client constructors
type Constructor func(url string, opts *config.Options) (Client, error)
