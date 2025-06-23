package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	log           *slog.Logger
	mux           *http.ServeMux
	server        *http.Server
	statusChannel chan PubSubMessage
	messages      []PubSubMessage
	mu            sync.Mutex
	sseConn       *SSEConn
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
		statusChannel: make(chan PubSubMessage),
	}
}

func (s *Server) Start() error {
	s.log.Info("Starting http server", "addr", s.server.Addr)

	s.setupRoutes()
	go s.processingIncoming()

	// Start the HTTP server
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	s.mux.HandleFunc("/messages", messageHandler(s.statusChannel, s.log))
	s.mux.HandleFunc("/sse-events", s.handleEvents())

	s.mux.HandleFunc("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return HomePage(s.messages), nil
	}))
}

func (s *Server) processingIncoming() {
	s.log.Info("Processing incoming messages...")
	for {
		msg := <-s.statusChannel
		s.log.Info("Received message", "message", msg)
		s.messages = append(s.messages, msg)
		s.broadcast("<div>test</div>", "incoming-messages")
	}
}

func (s *Server) handleEvents() http.HandlerFunc {
	s.sseConn = NewSSEConn(s.log)

	return func(w http.ResponseWriter, r *http.Request) {
		ch := s.sseConn.addClient(SSE_ALL_CLIENTS)
		defer s.sseConn.removeClient(SSE_ALL_CLIENTS, *ch)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			s.log.Debug("Could not init http.Flusher")
		}

		for {
			select {
			case message, ok := <-*ch:
				if ok {
					fmt.Println("case message... sending message")
					fmt.Println(message)
					_, _ = fmt.Fprintf(w, message)
					flusher.Flush()
				} else {
					return
				}
			case <-r.Context().Done():
				fmt.Println("Client closed connection")
				return
			}
		}
	}
}

func (s *Server) broadcast(data, event string) {
	s.sseConn.broadcast(SSE_ALL_CLIENTS, data, event)
}
