// Package flags contains description of all command-line flags
package flags

import (
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"
)

// InsecureSSLFLag allows to ignore certificate trust errors (eg for self-signed certificate)
var InsecureSSLFlag = altsrc.NewBoolFlag(&cli.BoolFlag{
	Name:    "insecure",
	Usage:   "Allow insecure server connections when using SSL",
	Value:   false,
	Aliases: []string{"k"},
})

// HeadersFlag allows to pass custom headers in websocket handshake request
var HeadersFlag = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
	Name:    "header",
	Usage:   "Pass custom header(s) to server. Also may be read from file with @file",
	Aliases: []string{"H"},
})

// PingPeriodFlag allows to configure period of sending "ping" frames
var PingPeriodFlag = altsrc.NewDurationFlag(&cli.DurationFlag{
	Name:    "ping-period",
	Usage:   "Send ping frames every <period>",
	Aliases: []string{"P"},
	Value:   0,
})

// IgnorePingsFlag allows to drop "ping" frames from server (not response with "pong" frames)
var IgnorePingsFlag = altsrc.NewBoolFlag(&cli.BoolFlag{
	Name:  "ignore-pings",
	Usage: "Do not send pong frames on received ping frames",
})

// TraceFlag allows to specify debug logs output
var TraceFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:  "trace",
	Usage: "Write a debug trace to FILE (\"-\" for STDOUT)",
})

// ShowHandshakeResponseFlag allows to save or print websocket handshake response
var ShowHandshakeResponseFlag = altsrc.NewBoolFlag(&cli.BoolFlag{
	Name:    "include",
	Aliases: []string{"i"},
	Usage:   "Include handshake response (headers and body) to output",
})

// SaveToConfigFlag allows to save current options to config file
var SaveConfigToFlag = &cli.StringFlag{
	Name:  "save-config",
	Usage: "Save current options to FILE (\"-\" for STDOUT)",
}

// ReadConfigFlag allows to read options from config
var ReadConfigFlag = &cli.StringFlag{
	Name:    "config",
	Aliases: []string{"K"},
	Usage:   "Read config from TOML FILE (\"-\" for STDIN)",
}

// OutputFlag allows to save received frames to file or print to stdout
var OutputFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:    "output",
	Aliases: []string{"o"},
	Usage:   "Write to FILE instead of STDOUT (\"-\" for STDOUT explicitly)",
})

// MessageAfterConnectFlag allows to send message to server after successful handshake
var MessageAfterConnectFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:    "data",
	Aliases: []string{"d"},
	Usage:   "Send a message to server after connection. Use @file notation to read from file or \"-\" to read from STDIN.",
})

// ProxyURLFlag allows to specify proxy address
var ProxyURLFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:    "proxy",
	Aliases: []string{"x"},
	Usage:   "[protocol://]host[:port] Use this proxy. Supported protocols: http, https, socks5",
})
