package util

import (
	"github.com/xakep666/wurl/pkg/client"
	"github.com/xakep666/wurl/pkg/client/gorilla"
	"github.com/xakep666/wurl/pkg/config"
	"gopkg.in/urfave/cli.v2"
)

const (
	optionsContextKey           = "options"
	clientConstructorContextKey = "client"
)

func SetupOptions(ctx *cli.Context) error {
	opts, err := OptionsFromFlags(ctx)
	if err != nil {
		return err
	}
	ctx.App.Metadata[optionsContextKey] = opts
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
