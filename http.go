package main

import (
	"context"
	"io"
	"log/slog"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
	"net/http"
	"time"
)

type Server struct {
	log    *slog.Logger
	mux    *http.ServeMux
	server *http.Server
}

type NewServerOptions struct {
	Log *slog.Logger
}

func NewServer(opts NewServerOptions) *Server {
	if opts.Log == nil {
		opts.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	mux := http.NewServeMux()

	return &Server{
		log: opts.Log,
		mux: mux,
		server: &http.Server{
			Addr:              ":8080",
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	s.log.Info("Starting http server", "addr", s.server.Addr)

	s.setupRoutes()

	// Start the HTTP server
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.log.Info("Stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("Error shutting down server", "error", err)
		return err
	}

	s.log.Info("HTTP server stopped gracefully")
	return nil
}

func (s *Server) setupRoutes() {
	Static(s.mux)

	s.mux.HandleFunc("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return HomePage(), nil
	}))
}
