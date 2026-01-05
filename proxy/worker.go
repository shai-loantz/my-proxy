package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
)

var headersBlacklist = map[string]struct{}{
	"Connection":          {},
	"Keep-Alive":          {},
	"Proxy-Authenticate":  {},
	"Proxy-Authorization": {},
	"TE":                  {},
	"Trailer":             {},
	"Transfer-Encoding":   {},
	"Upgrade":             {},
}

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
	server.state.total++
	if pr.ctx.Err() != nil {
		server.state.errors++
		return
	}

	log.Printf("New request is handeled in worker #%v\n", workerId)
	server.state.inflight++
	defer server.requestDone()

	upstreamReq, err := createUpstreamRequest(pr, *server.config.upstreamURL)
	if err != nil {
		http.Error(pr.responseWriter, "Bad Request", http.StatusBadRequest)
		server.state.errors++
		return
	}

	upstreamResponse, err := server.client.Do(upstreamReq)
	if err != nil {
		log.Println("Error while sending request to upstream", err)
		http.Error(pr.responseWriter, "Bad Gateway", http.StatusBadGateway)
		server.state.errors++
		return
	}

	err = copyResponse(upstreamResponse, pr.responseWriter)
	if err != nil {
		log.Println("Error while copying response to client", err)
		server.state.errors++
		return
	}
}

func createUpstreamRequest(pr *proxyRequest, upstreamURL url.URL) (*http.Request, error) {
	parsed, err := url.ParseRequestURI(pr.URLPath)
	if err != nil {
		log.Println("Error while parsing request URI", err)
		return nil, err
	}

	upstreamURL.Path = parsed.Path
	upstreamURL.RawQuery = parsed.RawQuery
	upstreamURL.Fragment = parsed.Fragment
	upstreamReq, err := http.NewRequestWithContext(
		pr.ctx,
		pr.method,
		upstreamURL.String(),
		bytes.NewReader(pr.body),
	)
	if err != nil {
		log.Println("Error while creating a new request object", err)
		return nil, err
	}
	for k, v := range pr.header {
		if _, ok := headersBlacklist[k]; ok {
			continue
		}
		upstreamReq.Header[k] = v
	}
	return upstreamReq, nil
}

func copyResponse(upstreamResponse *http.Response, responseWriter http.ResponseWriter) error {
	defer upstreamResponse.Body.Close()

	for k, vv := range upstreamResponse.Header {
		if _, ok := headersBlacklist[k]; ok {
			continue
		}
		for _, v := range vv {
			responseWriter.Header().Add(k, v)
		}
	}
	responseWriter.WriteHeader(upstreamResponse.StatusCode)

	_, err := io.Copy(responseWriter, upstreamResponse.Body)
	return err
}
