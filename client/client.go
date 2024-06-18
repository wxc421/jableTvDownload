package client

import (
	browser "github.com/EDDYCJY/fake-useragent"
	tlsclient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/go-resty/resty/v2"
	srt "github.com/juzeon/spoofed-round-tripper"
)

var proxyUrl string

func SetProxy(proxy string) {
	proxyUrl = proxy
}

func GetClient(httpClientOption ...tlsclient.HttpClientOption) (*resty.Client, error) {
	// Create a SpoofedRoundTripper that implements the http.RoundTripper interface
	tr, err := srt.NewSpoofedRoundTripper(
		httpClientOption...,
	)
	if err != nil {
		return nil, err
	}

	// Set as transport. Don't forget to set the UA!
	client := resty.New().SetTransport(tr).
		SetHeader("User-Agent", browser.Chrome())
	return client, nil
}

func GetProxyClient() (*resty.Client, error) {
	// Create a SpoofedRoundTripper that implements the http.RoundTripper interface
	tr, err := srt.NewSpoofedRoundTripper(
		tlsclient.WithRandomTLSExtensionOrder(), // needed for Chrome 107+
		tlsclient.WithProxyUrl(proxyUrl),
		tlsclient.WithClientProfile(profiles.Chrome_120),
	)

	if err != nil {
		return nil, err
	}

	if proxyUrl != "" {
		_ = tr.Client.SetProxy(proxyUrl) // clash etc.: socks5://127.0.0.1:7890
	}

	// Set as transport. Don't forget to set the UA!
	client := resty.New().SetTransport(tr).
		SetHeader("User-Agent", browser.Chrome())
	return client, nil
}
