package gorilla

import (
	"io"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"github.com/xakep666/wurl/pkg/client"
)

func (c *Client) Ping(payload []byte) error {
	c.log.Debugf("sending ping to %s", c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	return c.conn.WriteMessage(websocket.PingMessage, payload)
}

func (c *Client) WriteSingleMessage(payload []byte, messageType client.MessageType) error {
	c.log.Debugf("writing message (type %d) to %s", messageType, c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	return c.conn.WriteMessage(int(messageType), payload)
}

func (c *Client) WriteJSONMessage(obj interface{}) error {
	c.log.Debugf("writing json message to %s", c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	return c.conn.WriteJSON(obj)
}

const MessageSendBufferSize = 1024

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
