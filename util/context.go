package util

import (
	"io"
	"os"

	"github.com/xakep666/wurl/flags"
	"github.com/xakep666/wurl/pkg/client"
	"github.com/xakep666/wurl/pkg/client/gorilla"
	"github.com/xakep666/wurl/pkg/config"
	"gopkg.in/urfave/cli.v2"
)

const (
	optionsContextKey           = "options"
	clientConstructorContextKey = "client"
	outputContextKey            = "output"
)

func SetupOptions(ctx *cli.Context) error {
	var opts config.Options

	if ctx.IsSet(flags.ReadConfigFlag.Name) {
		if err := OptionsFromTOML(ctx, &opts); err != nil {
			return err
		}
	}

	err := OptionsFromFlags(ctx, &opts)
	if err != nil {
		return err
	}

	if ctx.IsSet(flags.SaveConfigToFlag.Name) {
		if err := OptionsToTOML(ctx, &opts); err != nil {
			return err
		}
	}

	ctx.App.Metadata[optionsContextKey] = &opts
	return nil
}

func MustGetOptions(ctx *cli.Context) *config.Options {
	opts, ok := ctx.App.Metadata[optionsContextKey]
	if !ok {
		panic("options not found in metadata")
	}
	return opts.(*config.Options)
}

func SetupClientConstructor(ctx *cli.Context) error {
	ctx.App.Metadata[clientConstructorContextKey] = client.Constructor(gorilla.NewClient)
	return nil
}

func MustGetClientConstructor(ctx *cli.Context) client.Constructor {
	cc, ok := ctx.App.Metadata[clientConstructorContextKey]
	if !ok {
		panic("client constructor not found in context")
	}
	return cc.(client.Constructor)
}

func SetupOutput(ctx *cli.Context) error {
	opts := MustGetOptions(ctx)
	switch opts.Output {
	case "", "-":
		ctx.App.Metadata[outputContextKey] = os.Stdout
		return nil
	default:
		file, err := os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			return err
		}
		ctx.App.Metadata[outputContextKey] = file
		return nil
	}
}

func MustGetOutput(ctx *cli.Context) io.WriteCloser {
	out, ok := ctx.App.Metadata[outputContextKey]
	if !ok {
		panic("output not found int context")
	}
	return out.(io.WriteCloser)
}
