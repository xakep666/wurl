package util

import (
	"github.com/urfave/cli/v2"
	"github.com/xakep666/wurl/flags"
	"github.com/xakep666/wurl/pkg/client"
	"github.com/xakep666/wurl/pkg/client/gorilla"
	"github.com/xakep666/wurl/pkg/config"
)

const (
	optionsContextKey           = "options"
	clientConstructorContextKey = "client"
)

// SetupOptions sets up application options.
// It calls "OptionsFromContext" to load options, optionally saves it to files using "OptionsToToml" and finally injects it to context.
func SetupOptions(ctx *cli.Context) error {
	opts, err := OptionsFromContext(ctx)
	if err != nil {
		return err
	}

	if ctx.IsSet(flags.SaveConfigToFlag.Name) {
		if err := OptionsToTOML(ctx); err != nil {
			return err
		}
	}

	ctx.App.Metadata[optionsContextKey] = opts
	return nil
}

// MustGetOptions extracts program options from context. It panics if options not found in context.
func MustGetOptions(ctx *cli.Context) *config.Options {
	opts, ok := ctx.App.Metadata[optionsContextKey]
	if !ok {
		panic("options not found in metadata")
	}
	return opts.(*config.Options)
}

// SetupClientConstructor injects client constructor to context.
func SetupClientConstructor(ctx *cli.Context) error {
	ctx.App.Metadata[clientConstructorContextKey] = client.Constructor(gorilla.NewClient)
	return nil
}

// MustGetClientConstructor extracts client constructor from context. It panics if it not found in context.
func MustGetClientConstructor(ctx *cli.Context) client.Constructor {
	cc, ok := ctx.App.Metadata[clientConstructorContextKey]
	if !ok {
		panic("client constructor not found in context")
	}
	return cc.(client.Constructor)
}
