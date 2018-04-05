package gorilla

import "github.com/gorilla/websocket"

func (c *Client) Ping(payload []byte) error {
	c.log.Printf("sending ping to %s", c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	return c.conn.WriteMessage(websocket.PingMessage, payload)
}
