package gorilla

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
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
	}
}

func NewClient(url string, opts *config.Options) (client.Client, error) {
	dialer := setupDialer(opts)
	conn, resp, err := dialer.Dial(url, opts.AdditionalHeaders)
	switch err {
	case nil:
		// pass
	case websocket.ErrBadHandshake:
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("bad handshake: can`t read body: %v", err)
		}
		return nil, fmt.Errorf("bad handshake: status %s, body %s", resp.Status, respBody)
	default:
		return nil, err
	}
	defer resp.Body.Close()

	log := logrus.StandardLogger().WithField("client", "gorilla")

	conn.SetPingHandler(func(appData string) error {
		if opts.RespondPings {
			log.Debugf("ping received from %s, payload %s", conn.RemoteAddr(), appData)
			return conn.WriteMessage(websocket.PingMessage, nil)
		}
		return nil
	})
	conn.SetPongHandler(func(appData string) error {
		log.Debugf("pong received from %s, payload %s", conn.RemoteAddr(), appData)
		return nil
	})

	if opts.PingPeriod > 0 {
		go func() {
			ticker := time.NewTicker(opts.PingPeriod)
			defer ticker.Stop()
			for {
				log.Debugf("sending ping message to %s", conn.RemoteAddr())
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.WithError(err).Error("websocket ping failed")
					return
				}

				<-ticker.C
			}
		}()
	}

	return &Client{conn: conn, opts: opts, log: log}, nil
}

func (c *Client) Close() error {
	c.log.Debugf("closing websocket client")
	return c.conn.Close()
}
