package tor

import (
	"fmt"
	"net/http"

	"golang.org/x/net/proxy"
)

const (
	// proxyAddr points to a local SOCKS proxy from Tor.
	proxyAddr = "127.0.0.1:9050"
)

// Client is an http.Client configured to connect to Tor.
type Client = http.Client

// NewClient returns an HTTP client that connect to a Tor proxy using SOCK5.
func NewClient() (*Client, error) {
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("error when creating Tor client")
	}

	httpTransport := &http.Transport{}
	httpTransport.Dial = dialer.Dial

	return &http.Client{Transport: httpTransport}, nil
}
