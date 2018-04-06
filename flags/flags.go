package flags

import (
	"gopkg.in/urfave/cli.v2"
)

var InsecureSSLFlag = cli.BoolFlag{
	Name:    "insecure",
	Usage:   "Allow insecure server connections when using SSL",
	Value:   false,
	Aliases: []string{"k"},
}

var HeadersFlag = cli.StringSliceFlag{
	Name:    "header",
	Usage:   "Pass custom header(s) to server. Also may be read from file with @file",
	Aliases: []string{"H"},
}

var PingPeriodFlag = cli.DurationFlag{
	Name:    "ping-period",
	Usage:   "Send ping frames every <period>",
	Aliases: []string{"P"},
	Value:   0,
}

var IgnorePingsFlag = cli.BoolFlag{
	Name:  "ignore-pings",
	Usage: "Do not send pong frames on received ping frames",
}

var TraceFlag = cli.StringFlag{
	Name:  "trace",
	Usage: "Write a debug trace to FILE (\"-\" for STDOUT)",
}

var ShowHandshakeResponseFlag = cli.BoolFlag{
	Name:    "include",
	Aliases: []string{"i"},
	Usage:   "Include handshake response (headers and body) to output",
}
