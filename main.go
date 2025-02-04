package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/flynnkc/home_db_microservice/pkg/handlers"
)

func main() {

	address := "127.0.0.1"
	port := "8080"
	tlsPort := "443"
	keyFile := ""
	certFile := ""
	log := slog.Default()

	m := GetMux()

	// Run server(s) in goroutine
	// Running HTTP & HTTPS concurrently is unreasonable due to the way hugo
	// writes URLs directly into pages, will not support both protocols
	var servers []*http.Server
	s := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", address, port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      m,
		//ErrorLog:     slog.NewLogLogger(handler, slog.LevelError),
	}
	servers = append(servers, s)
	go func() {
		log.Info("Starting HTTP Server...")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("error starting http server",
				"error", err.Error())
		}
	}()

	if keyFile != "" && certFile != "" {
		log.Info("Starting HTTPS Server..")
		s := &http.Server{
			Addr: fmt.Sprintf("%v:%v", address,
				tlsPort),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      m,
			//ErrorLog:     slog.NewLogLogger(eh, slog.LevelError),
			TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		}
		servers = append(servers, s)
		go func() {
			err := s.ListenAndServeTLS(certFile, keyFile)
			if err != nil && err != http.ErrServerClosed {
				log.Error("Server error", "Error", err)
			}
		}()
	}

	// Create channel that takes signal interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until signal
	<-c

	var wg sync.WaitGroup
	for _, server := range servers {

		wg.Add(1)
		go func(server *http.Server) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			defer cancel()
			defer wg.Done()

			log.Info("Shutting down server...", "Server", server.Addr)
			if err := server.Shutdown(ctx); err != nil {
				log.Error("Server shutdown failed", "error", err)
			}
			<-ctx.Done()
		}(server)
	}

	wg.Wait()
	log.Info("Shutdown Complete")
}

// GetMux initializes the multiplexer & handlers to de-clutter main
func GetMux() http.Handler {

	log := slog.Default()
	log.Debug("Setting handlers on Router")

	mux := http.NewServeMux()
	mux.Handle("/api/v1/read", http.HandlerFunc(handlers.MysqlReadHandler))

	// Middleware
	m := handlers.LoggingMiddleware(mux)
	m = handlers.RecoveryMiddleware(m)

	log.Debug("Returning new Mux",
		"Router", fmt.Sprintf("%+v", mux))

	return m
}
