package client

import (
	"errors"
	"os"
	"unicode/utf8"

	"github.com/xakep666/wurl/pkg/config"
)

var BinaryOutError = errors.New("binary output detected")

type BinaryCheckWriter struct {
	Opts *config.Options
}

func (b *BinaryCheckWriter) Write(p []byte) (n int, err error) {
	if !utf8.Valid(p) && b.Opts.Output == os.Stdout && !b.Opts.ForceBinaryToStdout {
		return 0, BinaryOutError
	}
	n, err = b.Opts.Output.Write(p)
	return
}
