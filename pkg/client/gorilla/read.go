package gorilla

import (
	"io"
)

// ReadTo reads data from connection and writes to provided writer.
func (c *Client) ReadTo(wr io.Writer) error {
	for {
		mtype, msg, err := c.conn.NextReader()
		if err != nil {
			return err
		}
		n, err := io.Copy(wr, msg)
		if err != nil {
			return err
		}
		c.log.Debugf("read %d bytes from message type %d", n, mtype)
	}
}
