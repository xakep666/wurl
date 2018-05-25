package client

import (
	"errors"
	"os"
	"unicode/utf8"

	"github.com/xakep666/wurl/pkg/config"
)

// BinaryOutError returned if client receives binary data and attempts to print if to stdout.
// However you can still print binary data if "output" flag explicitly set to "-".
var BinaryOutError = errors.New("binary output detected")

// BinaryCheckWriter checks if we attempt to write binary data to terminal.
type BinaryCheckWriter struct {
	Opts *config.Options
}

// Writer implements io.Writer interface.
func (b *BinaryCheckWriter) Write(p []byte) (n int, err error) {
	if !utf8.Valid(p) && b.Opts.Output == os.Stdout && !b.Opts.ForceBinaryToStdout {
		return 0, BinaryOutError
	}
	n, err = b.Opts.Output.Write(p)
	return
}
