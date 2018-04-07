package commands

import (
	"fmt"

	"github.com/xakep666/wurl/flags"
	"github.com/xakep666/wurl/pkg/client"
	"github.com/xakep666/wurl/util"
	"gopkg.in/urfave/cli.v2"
)

var ReadCommand = cli.Command{
	Name:      "read",
	Usage:     "Simply read from websocket",
	ArgsUsage: "<url>",
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() < 1 {
			cli.ShowCommandHelp(ctx, "read")
			return fmt.Errorf("url must be provided")
		}
		cl, err := util.MustGetClientConstructor(ctx)(ctx.Args().First(), util.MustGetOptions(ctx))
		if err != nil {
			return err
		}
		defer cl.Close()

		opts := util.MustGetOptions(ctx)
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
	},
}
