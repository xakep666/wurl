package util

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/xakep666/wurl/flags"
	"github.com/xakep666/wurl/pkg/config"
	"golang.org/x/net/proxy"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"
)

func processHeadersFlag(values []string) (ret http.Header, err error) {
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

type nopWriteCloser struct {
	io.WriteCloser
}

func (n *nopWriteCloser) Close() error {
	return nil
}

func processOutToFlag(outOpt string) (io.WriteCloser, error) {
	switch outOpt {
	case "", "-":
		return &nopWriteCloser{os.Stdout}, nil
	default:
		return os.OpenFile(outOpt, os.O_APPEND|os.O_CREATE, os.ModePerm)
	}
}

func processTraceToFlag(outOpt string) (io.WriteCloser, error) {
	switch outOpt {
	case "":
		return nil, nil
	case "-":
		return &nopWriteCloser{os.Stdout}, nil
	default:
		return os.OpenFile(outOpt, os.O_APPEND|os.O_CREATE, os.ModePerm)
	}
}

func processFromFlag(inOpt string) (io.ReadCloser, error) {
	switch {
	case inOpt == "":
		return nil, nil
	case inOpt == "-":
		return os.Stdin, nil
	case strings.HasPrefix(inOpt, "@"):
		return os.Open(strings.TrimPrefix(inOpt, "@"))
	default:
		return ioutil.NopCloser(strings.NewReader(inOpt)), nil
	}
}

func processProxyFlag(urlOpt string) (dialFunc config.DialFunc, err error) {
	proxyURL, err := url.Parse(urlOpt)
	if err != nil {
		return nil, err
	}

	switch proxyURL.Scheme {
	case "http", "https":
		var dialer proxy.Dialer
		dialer, err = proxy.FromURL(proxyURL, proxy.Direct)
		dialFunc = dialer.Dial
	case "socks5":
		var auth *proxy.Auth
		if proxyURL.User != nil {
			auth = &proxy.Auth{User: proxyURL.User.Username()}
			auth.Password, _ = proxyURL.User.Password()
		}
		var dialer proxy.Dialer
		dialer, err = proxy.SOCKS5("tcp", proxyURL.Host, auth, proxy.Direct)
		dialFunc = dialer.Dial
	default:
		err = fmt.Errorf("unsupported proxy protocol \"%s\"", proxyURL.Scheme)
	}

	return
}

func OptionsFromContext(ctx *cli.Context) (opts *config.Options, err error) {
	opts = &config.Options{}
	opts.AllowInsecureSSL = ctx.Bool(flags.InsecureSSLFlag.Name)
	opts.PingPeriod = ctx.Duration(flags.PingPeriodFlag.Name)
	opts.RespondPings = !ctx.Bool(flags.IgnorePingsFlag.Name)
	opts.AdditionalHeaders, err = processHeadersFlag(ctx.StringSlice(flags.HeadersFlag.Name))
	if err != nil {
		return
	}
	opts.TraceTo, err = processTraceToFlag(ctx.String(flags.TraceFlag.Name))
	if err != nil {
		return
	}
	opts.ShowHandshakeResponse = ctx.Bool(flags.ShowHandshakeResponseFlag.Name)
	opts.Output, err = processOutToFlag(ctx.String(flags.OutputFlag.Name))
	if err != nil {
		return
	}
	opts.ForceBinaryToStdout = ctx.IsSet(flags.OutputFlag.Name)
	opts.MessageAfterConnect, err = processFromFlag(ctx.String(flags.MessageAfterConnectFlag.Name))
	if err != nil {
		return
	}
	opts.DialFunc, err = processProxyFlag(ctx.String(flags.ProxyURLFlag.Name))

	return
}

func OptionsToTOML(ctx *cli.Context) error {
	fileName := ctx.String(flags.SaveConfigToFlag.Name)
	var out io.Writer
	switch fileName {
	case "":
		return nil
	case "-":
		out = os.Stdout
	default:
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer file.Close()
		out = file
	}

	optMap := make(map[string]interface{})
	for _, option := range ctx.App.Flags {
		optName := option.Names()[0]
		switch option.(type) {
		// encode only "altsrc" flags
		case *altsrc.BoolFlag,
			*altsrc.DurationFlag,
			*altsrc.Float64Flag,
			*altsrc.GenericFlag,
			*altsrc.Int64Flag,
			*altsrc.IntFlag,
			*altsrc.PathFlag,
			*altsrc.StringFlag,
			*altsrc.Uint64Flag,
			*altsrc.UintFlag:
			optMap[optName] = ctx.Generic(optName)
		case *altsrc.Float64SliceFlag:
			optMap[optName] = ctx.Float64Slice(optName)
		case *altsrc.Int64SliceFlag:
			optMap[optName] = ctx.Int64Slice(optName)
		case *altsrc.IntSliceFlag:
			optMap[optName] = ctx.IntSlice(optName)
		case *altsrc.StringSliceFlag:
			optMap[optName] = ctx.StringSlice(optName)
		}
	}
	fmt.Println(optMap)
	return toml.NewEncoder(out).Encode(optMap)
}

func SetupLogger(ctx *cli.Context) {
	traceTo := MustGetOptions(ctx).TraceTo
	if traceTo == nil {
		logrus.SetLevel(logrus.ErrorLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(traceTo)
	}
}
