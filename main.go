package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/blang/semver"
	"github.com/sirupsen/logrus"
	"github.com/xakep666/wurl/commands"
	"github.com/xakep666/wurl/flags"
	"github.com/xakep666/wurl/util"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"
)

var Version = semver.MustParse("0.0.1-alpha")

func main() {
	app := cli.App{
		Name:    "wurl",
		Usage:   "console websocket client",
		Version: Version.String(),
		Flags: []cli.Flag{
			flags.InsecureSSLFlag,
			flags.HeadersFlag,
			flags.PingPeriodFlag,
			flags.IgnorePingsFlag,
			flags.TraceFlag,
			flags.ShowHandshakeResponseFlag,
			flags.ReadConfigFlag,
			flags.SaveConfigToFlag,
			flags.OutputFlag,
			flags.MessageAfterConnectFlag,
		},
		Commands: []*cli.Command{
			&commands.ReadCommand,
		},
	}

	app.Before = func(ctx *cli.Context) error {
		if ctx.IsSet(flags.ReadConfigFlag.Name) {
			loadFromConfig := altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewTomlSourceFromFlagFunc(flags.ReadConfigFlag.Name))
			if err := loadFromConfig(ctx); err != nil {
				return err
			}
		}
		if err := setup(ctx); err != nil {
			return err
		}
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		opts := reflect.ValueOf(util.MustGetOptions(ctx))
		for i := 0; i < opts.NumField(); i++ {
			if closer, ok := opts.Field(i).Interface().(io.Closer); ok {
				closer.Close()
			}
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func setup(ctx *cli.Context) error {
	if err := util.SetupOptions(ctx); err != nil {
		return err
	}
	util.SetupLogger(ctx)
	logrus.Debugf("running with config %+v", util.MustGetOptions(ctx))
	if err := util.SetupClientConstructor(ctx); err != nil {
		return err
	}

	return nil
}
