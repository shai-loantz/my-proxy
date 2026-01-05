package main

import "flag"

type args struct {
	listenAddress string
	upstreamURL   string
	workersNum    uint
	queueSize     uint
}

func (a *args) Register() {
	flag.StringVar(&a.listenAddress, "listen", "0.0.0.0:8080", "address to listen to")
	flag.StringVar(&a.upstreamURL, "upstreamURL", "", "Where should the proxy forward all the requests")
	flag.UintVar(&a.workersNum, "workers", 10, "Number of workers")
	flag.UintVar(&a.queueSize, "queue-size", 50, "Size of the queue")
}

func getArgs() args {
	var a args
	a.Register()
	flag.Parse()
	return a
}
