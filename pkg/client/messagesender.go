package client

import (
	"io"
	"unicode/utf8"
)

const MessageSendBufferSize = 1024

type MessageSender struct {
	Reader io.Reader
}

func (m *MessageSender) SendMessage(client Client) error {
	buf := [MessageSendBufferSize]byte{}
	for {
		n, err := m.Reader.Read(buf[:])
		switch err {
		case nil:
			// pass
		case io.EOF:
			return nil
		default:
			return err
		}
		var messageType MessageType
		if utf8.Valid(buf[:n]) {
			messageType = TextMessage
		} else {
			messageType = BinaryMessage
		}

		if err := client.WriteSingleMessage(buf[:n], messageType); err != nil {
			return err
		}
	}
}
