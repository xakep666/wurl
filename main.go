package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/blang/semver"
	"github.com/sirupsen/logrus"
	"github.com/xakep666/wurl/flags"
	"github.com/xakep666/wurl/pkg/client"
	"github.com/xakep666/wurl/util"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"
)

var Version = semver.MustParse("0.0.1-alpha")

func main() {
	app := &cli.App{
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
			flags.ProxyURLFlag,
			// completion
			util.InitCompletionFlag,
		},
		CustomAppHelpTemplate: util.AppHelp,
		Commands:              nil,
		Before:                setup,
		Action:                action,
		After:                 after,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func setup(ctx *cli.Context) error {
	if ctx.IsSet(util.InitCompletionFlag.Name) {
		util.PrintCompletion(ctx)
		return cli.Exit("", 0)
	}

	if ctx.IsSet(flags.ReadConfigFlag.Name) {
		loadFromConfig := altsrc.InitInputSourceWithContext(ctx.App.Flags, altsrc.NewTomlSourceFromFlagFunc(flags.ReadConfigFlag.Name))
		if err := loadFromConfig(ctx); err != nil {
			return cli.Exit(err, 1)
		}
	}

	if err := util.SetupOptions(ctx); err != nil {
		return cli.Exit(err, 1)
	}

	util.SetupLogger(ctx)
	logrus.Debugf("running with config %+v", util.MustGetOptions(ctx))
	if err := util.SetupClientConstructor(ctx); err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func action(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		cli.ShowAppHelp(ctx)
		return fmt.Errorf("url must be provided")
	}

	opts := util.MustGetOptions(ctx)

	cl, resp, err := util.MustGetClientConstructor(ctx)(ctx.Args().First(), opts)
	if opts.ShowHandshakeResponse {
		if resp != nil {
			resp.Write(os.Stdout)
			fmt.Fprintln(os.Stdout)
		}
	}
	if err != nil {
		return err
	}
	defer cl.Close()

	if opts.MessageAfterConnect != nil {
		if err := cl.WriteMessageFrom(opts.MessageAfterConnect); err != nil {
			return err
		}
	}

	err = cl.ReadTo(&client.BinaryCheckWriter{Opts: opts})
	switch err {
	case nil:
		// pass
	case client.BinaryOutError:
		fmt.Println("WARNING: binary output can mess up your terminal.")
		fmt.Printf("Use \"--%[1]s -\" to tell %[2]s to output it to your "+
			"terminal anyway, or consider \"--%[1]s <FILE>\" to save to a file.\n", flags.OutputFlag.Name, ctx.App.Name)
		return nil
	default:
		return err
	}

	return nil
}

func after(ctx *cli.Context) error {
	opts := reflect.ValueOf(util.MustGetOptions(ctx)).Elem()
	for i := 0; i < opts.NumField(); i++ {
		if closer, ok := opts.Field(i).Interface().(io.Closer); ok {
			closer.Close()
		}
	}
	return nil
}
