// Package gorilla contains websocket client implementation using "github.com/gorilla/websocket" library.
package gorilla

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/xakep666/wurl/pkg/client"
	"github.com/xakep666/wurl/pkg/config"
)

var (
	_ client.Client      = &Client{}
	_ client.Constructor = NewClient
)

// Client is a websocket client implementation
type Client struct {
	conn           *websocket.Conn
	connWriteMutex sync.Mutex

	opts *config.Options
	log  *logrus.Entry
}

func setupDialer(opts *config.Options) *websocket.Dialer {
	return &websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.AllowInsecureSSL,
		},
		NetDial: opts.DialFunc,
	}
}

func (c *Client) prepareConnection() {
	c.conn.SetPingHandler(func(appData string) error {
		if c.opts.RespondPings {
			c.log.Debugf("ping received from %s, payload %s", c.conn.RemoteAddr(), appData)
			return c.conn.WriteMessage(websocket.PongMessage, nil)
		}
		return nil
	})
	c.conn.SetPongHandler(func(appData string) error {
		c.log.Debugf("pong received from %s, payload %s", c.conn.RemoteAddr(), appData)
		return nil
	})
}

func (c *Client) periodicPinger() {
	ticker := time.NewTicker(c.opts.PingPeriod)
	defer ticker.Stop()
	for {
		if err := c.Ping(nil); err != nil {
			c.log.WithError(err).Error("websocket ping failed")
			return
		}

		<-ticker.C
	}
}

// NewClient constructs a websocket client.
// First stage is a connection establishing.
// Second stage is connection tuning: adding ping/pong handlers according to options.
// Third stage is a optional start of periodic ping routine (if specified).
func NewClient(url string, opts *config.Options) (client.Client, *http.Response, error) {
	dialer := setupDialer(opts)
	conn, resp, err := dialer.Dial(url, opts.AdditionalHeaders)

	switch err {
	case nil:
		// pass
	case websocket.ErrBadHandshake:
		return nil, resp, client.ErrBadHandshake
	default:
		return nil, resp, err
	}

	log := logrus.StandardLogger().WithField("client", "gorilla")

	ret := &Client{conn: conn, opts: opts, log: log}

	ret.prepareConnection()

	if opts.PingPeriod > 0 {
		go ret.periodicPinger()
	}

	return ret, resp, nil
}

// Close closes connection to server.
func (c *Client) Close() error {
	c.log.Debugf("closing websocket client")
	return c.conn.Close()
}
