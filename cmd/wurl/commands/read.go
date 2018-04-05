package commands

import (
	"fmt"
	"os"

	"github.com/xakep666/wurl/cmd/wurl/util"
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

		if err := cl.ReadTo(os.Stdout); err != nil {
			return err
		}

		return nil
	},
}
