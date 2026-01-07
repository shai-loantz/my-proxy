package proxy

import "net/url"

type serverConfig struct {
	upstreamURL   *url.URL
	listenAddress string
}

func NewServerConfig(upstreamURL string, listenAddress string) (serverConfig, error) {
	parsedURL, err := url.Parse(upstreamURL)
	if err != nil {
		return serverConfig{}, err
	}

	return serverConfig{
		upstreamURL:   parsedURL,
		listenAddress: listenAddress,
	}, nil
}
