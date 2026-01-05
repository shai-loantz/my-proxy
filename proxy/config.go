package proxy

import "net/url"

type serverConfig struct {
	upstreamURL *url.URL
	queueSize   uint
	workersNum  uint
}

func NewServerConfig(upstreamURL string, queueSize uint, workerNum uint) (serverConfig, error) {
	parsedURL, err := url.Parse(upstreamURL)
	if err != nil {
		return serverConfig{}, err
	}

	return serverConfig{
		upstreamURL: parsedURL,
		queueSize:   queueSize,
		workersNum:  workerNum,
	}, nil
}
