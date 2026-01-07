package main

import (
	"log"
	"my-proxy/proxy"
	"net/http"
)

func main() {
	args := getArgs()
	config, err := proxy.NewServerConfig(
		args.upstreamURL,
		args.queueSize,
		args.workersNum,
	)
	if err != nil {
		log.Fatalln("Could not create server config.", err)
	}

	log.Printf("Starting proxy server on %s, upstream=%s, workersNum=%d, queue-size=%d\n", args.listenAddress, args.upstreamURL, args.workersNum, args.queueSize)
	proxyServer, err := proxy.NewServer(config)
	log.Println("Proxy server created. Serving...")
	err = http.ListenAndServe(args.listenAddress, proxyServer)
	if err != nil {
		proxyServer.Shutdown()
		log.Fatalln(err)
	}
}
