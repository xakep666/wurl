package client

import (
	"errors"
	"io"
	"os"
	"unicode/utf8"

	"github.com/xakep666/wurl/pkg/config"
)

var BinaryOutError = errors.New("binary output detected")

type BinaryCheckWriter struct {
	Opts *config.Options
	io.Writer
}

func (b *BinaryCheckWriter) Write(p []byte) (n int, err error) {
	if !utf8.Valid(p) && b.Writer == os.Stdout && b.Opts.Output == "" {
		return 0, BinaryOutError
	}
	n, err = b.Writer.Write(p)
	return
}
