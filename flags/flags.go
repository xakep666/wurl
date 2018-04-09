package flags

import (
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"
)

var InsecureSSLFlag = altsrc.NewBoolFlag(&cli.BoolFlag{
	Name:    "insecure",
	Usage:   "Allow insecure server connections when using SSL",
	Value:   false,
	Aliases: []string{"k"},
})

var HeadersFlag = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
	Name:    "header",
	Usage:   "Pass custom header(s) to server. Also may be read from file with @file",
	Aliases: []string{"H"},
})

var PingPeriodFlag = altsrc.NewDurationFlag(&cli.DurationFlag{
	Name:    "ping-period",
	Usage:   "Send ping frames every <period>",
	Aliases: []string{"P"},
	Value:   0,
})

var IgnorePingsFlag = altsrc.NewBoolFlag(&cli.BoolFlag{
	Name:  "ignore-pings",
	Usage: "Do not send pong frames on received ping frames",
})

var TraceFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:  "trace",
	Usage: "Write a debug trace to FILE (\"-\" for STDOUT)",
})

var ShowHandshakeResponseFlag = altsrc.NewBoolFlag(&cli.BoolFlag{
	Name:    "include",
	Aliases: []string{"i"},
	Usage:   "Include handshake response (headers and body) to output",
})

var SaveConfigToFlag = &cli.StringFlag{
	Name:  "save-config",
	Usage: "Save current options to FILE (\"-\" for STDOUT)",
}

var ReadConfigFlag = &cli.StringFlag{
	Name:    "config",
	Aliases: []string{"K"},
	Usage:   "Read config from TOML FILE (\"-\" for STDIN)",
}

var OutputFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:    "output",
	Aliases: []string{"o"},
	Usage:   "Write to FILE instead of STDOUT (\"-\" for STDOUT explicitly)",
})

var MessageAfterConnectFlag = altsrc.NewStringFlag(&cli.StringFlag{
	Name:    "data",
	Aliases: []string{"d"},
	Usage:   "Send a message to server after connection. Use @file notation to read from file or \"-\" to read from STDIN.",
})
