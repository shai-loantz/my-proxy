package proxy

import (
	"log"
	"net/http"
	"time"
)

type serverState struct {
	inflight int64
	total    int64
	errors   int64
}

type Server struct {
	config     serverConfig
	state      *serverState
	client     *http.Client
	HttpServer *http.Server
}

func NewServer(config serverConfig) (*Server, error) {
	state := &serverState{0, 0, 0}
	client := &http.Client{Timeout: 60 * time.Second}
	httpServer := &http.Server{
		Addr:    config.listenAddress,
		Handler: nil,
	}

	server := &Server{
		config:     config,
		state:      state,
		client:     client,
		HttpServer: httpServer,
	}

	return server, nil
}

func (server *Server) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Got a new request. Parsing")

	proxyReq := NewProxyRequest(responseWriter, request, request.Context())
	server.handleProxyRequest(&proxyReq)
}

func (server *Server) requestDone() {
	server.state.inflight--
}
