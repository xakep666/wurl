package commands

import (
	"fmt"
	"io"

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

		msg := util.GetSingleMessageReader(ctx)
		if msg != nil {
			defer msg.Close()
			if err := cl.WriteMessageFrom(msg); err != nil {
				return err
			}
		}

		var out io.Writer = &client.BinaryCheckWriter{
			Writer: util.MustGetOutput(ctx),
			Opts:   util.MustGetOptions(ctx),
		}

		err = cl.ReadTo(out)
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
