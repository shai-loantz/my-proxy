package proxy

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
)

func (server *Server) worker(workerId int) {
	defer server.wg.Done()

	for {
		select {
		case pr := <-server.queue:
			server.handleProxyRequest(pr, workerId)
		case <-server.ctx.Done():
			return
		}
	}
}

func (server *Server) handleProxyRequest(pr *proxyRequest, workerId int) {
	if pr.ctx.Err() != nil {
		server.state.errors++
		return
	}

	log.Printf("New request is handeled in worker #%v\n", workerId)
	server.state.inflight++
	defer server.requestDone()
	server.state.total++

	newRequestURL := *server.config.upstreamURL

	parsed, err := url.ParseRequestURI(pr.URLPath)
	if err != nil {
		log.Println("Error while parsing request URI", err)
		http.Error(pr.responseWriter, "Bad Request URI", http.StatusBadRequest)
		server.state.errors++
		return
	}

	newRequestURL.Path = parsed.Path
	newRequestURL.RawQuery = parsed.RawQuery
	newRequestURL.Fragment = parsed.Fragment
	newRequest, err := http.NewRequestWithContext(
		pr.ctx,
		pr.method,
		newRequestURL.String(),
		bytes.NewReader(pr.body),
	)
	// TODO: copy headers

	server.client.Do(newRequest)
	// TODO: handle errors and return response to original client
}
