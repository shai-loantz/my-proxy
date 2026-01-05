package proxy

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type serverState struct {
	inflight int64
	total    int64
	errors   int64
}

type Server struct {
	config serverConfig
	state  *serverState
	queue  chan *proxyRequest
	client *http.Client

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func NewServer(config serverConfig) (*Server, error) {
	state := &serverState{0, 0, 0}
	client := &http.Client{Timeout: 60 * time.Second}
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	queue := make(chan *proxyRequest, config.queueSize)

	server := &Server{
		config: config,
		state:  state,
		queue:  queue,
		client: client,
		ctx:    ctx,
		cancel: cancel,
		wg:     wg,
	}

	server.wg.Add(int(config.workersNum))
	for i := 0; i < int(config.workersNum); i++ {
		go server.worker(i)
	}

	return server, nil
}

func (server *Server) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Got a new request. Parsing")
	proxyReq := NewProxyRequest(responseWriter, request)
	select {
	case server.queue <- &proxyReq:
		log.Println("Queued")
	case <-time.After(time.Second):
		log.Println("Timeout while queuing (queue is full)")
		http.Error(responseWriter, "Proxy Timeout", http.StatusGatewayTimeout)
	}
}

func (server *Server) Shutdown() {
	server.cancel()
	server.wg.Wait()
}

func (server *Server) requestDone() {
	server.state.inflight--
}
