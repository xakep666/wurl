package main

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/sirupsen/logrus"
	"github.com/xakep666/wurl/cmd/wurl/commands"
	"github.com/xakep666/wurl/cmd/wurl/flags"
	"github.com/xakep666/wurl/cmd/wurl/util"
	"gopkg.in/urfave/cli.v2"
)

var Version = semver.MustParse("0.0.1-alpha")

func main() {
	app := cli.App{
		Name:    "wurl",
		Usage:   "console websocket client",
		Version: Version.String(),
		Flags: []cli.Flag{
			&flags.InsecureSSLFlag,
			&flags.HeadersFlag,
			&flags.PingPeriodFlag,
			&flags.IgnorePingsFlag,
			&flags.TraceFlag,
		},
		Commands: []*cli.Command{
			&commands.ReadCommand,
		},
		Before: func(ctx *cli.Context) error {
			if err := util.SetupOptions(ctx); err != nil {
				return err
			}
			if err := util.SetupLogger(ctx); err != nil {
				return err
			}
			logrus.Debugf("running with config %+v", util.MustGetOptions(ctx))
			if err := util.SetupClientConstructor(ctx); err != nil {
				return err
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
