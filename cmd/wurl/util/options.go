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
	"github.com/xakep666/wurl/cmd/wurl/flags"
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

func OptionsFromFlags(ctx *cli.Context) (ret *config.Options, err error) {
	ret = &config.Options{}
	ret.AllowInsecureSSL = ctx.Bool(flags.InsecureSSLFlag.Name)
	ret.PingPeriod = ctx.Duration(flags.PingPeriodFlag.Name)
	ret.RespondPings = !ctx.Bool(flags.IgnorePingsFlag.Name)
	ret.AdditionalHeaders, err = processHeadersFlag(ctx)
	if err != nil {
		return
	}
	ret.TraceTo = ctx.String(flags.TraceFlag.Name)

	return
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
