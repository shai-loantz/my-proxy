package main

import (
	"log"
	"my-proxy/proxy"
	"net/http"
)

func main() {
	args := getArgs()
	log.Printf("listen=%v, upstream=%v, workers=%v, queue-size=%v\n", args.listenAddress, args.upstreamURL, args.workersNum, args.queueSize)
	config, err := proxy.NewServerConfig(
		args.upstreamURL,
		args.queueSize,
		args.workersNum,
	)
	if err != nil {
		log.Fatalln("Could not create server config.", err)
	}

	proxyServer, err := proxy.NewServer(config)
	log.Fatalln(http.ListenAndServe(args.listenAddress, proxyServer))
}
