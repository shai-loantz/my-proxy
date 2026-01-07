package main

import (
	"context"
	"errors"
	"log"
	"my-proxy/proxy"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	args := getArgs()
	config, err := proxy.NewServerConfig(
		args.upstreamURL,
		args.listenAddress,
	)
	if err != nil {
		log.Fatalln("Could not create server config.", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("Starting proxy server on %s, upstream=%s\n", args.listenAddress, args.upstreamURL)
	proxyServer, err := proxy.NewServer(config)
	http.Handle("/", proxyServer)
	log.Println("Proxy server created. Serving...")

	go func() {
		if err = proxyServer.HttpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down proxy server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := proxyServer.HttpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
