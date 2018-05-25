package gorilla

import (
	"io"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"github.com/xakep666/wurl/pkg/client"
)

// Ping sends ping frame to server with optionally provided payload.
func (c *Client) Ping(payload []byte) error {
	c.log.Debugf("sending ping to %s", c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	return c.conn.WriteMessage(websocket.PingMessage, payload)
}

// WriteSingleMessage sends given payload to server. Frame type set to "messageType".
func (c *Client) WriteSingleMessage(payload []byte, messageType client.MessageType) error {
	c.log.Debugf("writing message (type %d) to %s", messageType, c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	return c.conn.WriteMessage(int(messageType), payload)
}

// WriteJSONMessage sends json-encoded object to server.
func (c *Client) WriteJSONMessage(obj interface{}) error {
	c.log.Debugf("writing json message to %s", c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	return c.conn.WriteJSON(obj)
}

// MessageSendBufferSize is a send buffer size (in bytes)
const MessageSendBufferSize = 1024

// WriteMessageFrom sends bytes read from provided reader until it returns io.EOF
func (c *Client) WriteMessageFrom(reader io.Reader) error {
	c.log.Debugf("writing message from reader to %s", c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	buf := [MessageSendBufferSize]byte{}
	for {
		n, err := reader.Read(buf[:])
		switch err {
		case nil:
			// pass
		case io.EOF:
			return nil
		default:
			return err
		}

		var messageType int
		if utf8.Valid(buf[:n]) {
			messageType = websocket.TextMessage
		} else {
			messageType = websocket.BinaryMessage
		}

		writer, err := c.conn.NextWriter(messageType)
		if err != nil {
			return err
		}

		if _, err := writer.Write(buf[:n]); err != nil {
			return err
		}
	}
}
