package gorilla

import "github.com/gorilla/websocket"

func (c *Client) Ping(payload []byte) error {
	c.log.Debugf("sending ping to %s", c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	return c.conn.WriteMessage(websocket.PingMessage, payload)
}

func (c *Client) WriteSingleMessage(payload []byte, messageType int) error {
	c.log.Debugf("writing message (type %d) to %s", messageType, c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	return c.conn.WriteMessage(messageType, payload)
}

func (c *Client) WriteJSONMessage(obj interface{}) error {
	c.log.Debugf("writing json message to %s", c.conn.RemoteAddr())
	c.connWriteMutex.Lock()
	defer c.connWriteMutex.Unlock()
	return c.conn.WriteJSON(obj)
}
