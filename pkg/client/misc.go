package client

import (
	"fmt"
	"net/http"
	"os"

	"github.com/xakep666/wurl/pkg/config"
)

func WriteHandshakeResponse(resp *http.Response, opts *config.Options) error {
	if opts.ShowHandshakeResponse {
		if err := resp.Write(os.Stdout); err != nil {
			return err
		}

		_, err := fmt.Fprintln(os.Stdout)
		return err
	}
	return nil
}
