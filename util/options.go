package util

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xakep666/wurl/flags"
	"github.com/xakep666/wurl/pkg/config"
	"gopkg.in/urfave/cli.v2"
)

func processHeadersFlag(ctx *cli.Context) (ret http.Header, err error) {
	values := ctx.StringSlice(flags.HeadersFlag.Name)
	ret = make(http.Header)

	for _, value := range values {
		var rc io.ReadCloser
		if strings.HasPrefix(value, "@") {
			rc, err = os.Open(strings.TrimPrefix(value, "@"))
			if err != nil {
				return
			}
		} else {
			rc = ioutil.NopCloser(strings.NewReader(value))
		}

		var mimeHeader textproto.MIMEHeader
		mimeHeader, err = textproto.NewReader(bufio.NewReader(rc)).ReadMIMEHeader()
		switch err {
		case nil, io.EOF:
			err = nil
		default:
			rc.Close()
			return
		}

		for name, values := range mimeHeader {
			for _, value := range values {
				ret.Add(name, value)
			}
		}

		rc.Close()
	}
	return
}

func OptionsFromFlags(ctx *cli.Context, opts *config.Options) (err error) {
	if ctx.IsSet(flags.InsecureSSLFlag.Name) {
		opts.AllowInsecureSSL = ctx.Bool(flags.InsecureSSLFlag.Name)
	}

	if ctx.IsSet(flags.PingPeriodFlag.Name) {
		opts.PingPeriod = ctx.Duration(flags.PingPeriodFlag.Name)
	}

	if ctx.IsSet(flags.IgnorePingsFlag.Name) {
		opts.RespondPings = !ctx.Bool(flags.IgnorePingsFlag.Name)
	}

	if ctx.IsSet(flags.HeadersFlag.Name) {
		opts.AdditionalHeaders, err = processHeadersFlag(ctx)
		if err != nil {
			return
		}
	}

	if ctx.IsSet(flags.TraceFlag.Name) {
		opts.TraceTo = ctx.String(flags.TraceFlag.Name)
	}

	if ctx.IsSet(flags.ShowHandshakeResponseFlag.Name) {
		opts.ShowHandshakeResponse = ctx.Bool(flags.ShowHandshakeResponseFlag.Name)
	}

	return
}

func OptionsFromTOML(ctx *cli.Context, opts *config.Options) (err error) {
	fileName := ctx.String(flags.ReadConfigFlag.Name)
	switch fileName {
	case "":
		return
	case "-":
		return config.FromTOML(os.Stdin, opts)
	default:
		var file *os.File
		file, err = os.Open(fileName)
		if err != nil {
			return
		}
		defer file.Close()
		err = config.FromTOML(file, opts)
		return
	}
}

func OptionsToTOML(ctx *cli.Context, opts *config.Options) (err error) {
	fileName := ctx.String(flags.SaveConfigToFlag.Name)
	switch fileName {
	case "":
		return
	case "-":
		return opts.ToTOML(os.Stdout)
	default:
		var file *os.File
		file, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return
		}
		defer file.Close()
		err = opts.ToTOML(file)
		return
	}
}

func SetupLogger(ctx *cli.Context) error {
	traceTo := MustGetOptions(ctx).TraceTo
	switch traceTo {
	case "":
		logrus.SetLevel(logrus.ErrorLevel)
	case "-":
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stdout)
	default:
		file, err := os.OpenFile(traceTo, os.O_APPEND, os.ModePerm)
		if err != nil {
			return err
		}
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(file)
	}
	return nil
}
