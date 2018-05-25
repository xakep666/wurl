package client

import (
	"errors"
)

// ErrBadHandshake returned by client constructor on bad handshake
var ErrBadHandshake = errors.New("bad handshake")
