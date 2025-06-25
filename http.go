package main

import (
	"bytes"
	"errors"
	"fmt"

	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/r3labs/sse/v2"

	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

type Server struct {
	opts          NewServerOptions
	log           *slog.Logger
	router        *chi.Mux
	statusChannel chan PubSubMessage
	messages      []PubSubMessage
	mu            sync.Mutex
	sse           *sse.Server
}

type NewServerOptions struct {
	Log         *slog.Logger
	Addr        string
	OfflineMode bool
}

func NewServer(opts NewServerOptions) *Server {
	if opts.Log == nil {
		opts.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	if opts.Addr == "" {
		opts.Addr = ":8080"
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(opts.Log, &httplog.Options{Level: slog.LevelInfo}))
	r.Use(middleware.Recoverer)

	return &Server{
		log:           opts.Log,
		router:        r,
		statusChannel: make(chan PubSubMessage),
		opts:          opts,
	}
}

func (s *Server) Start() error {
	s.log.Info("Starting http server", "addr", s.opts.Addr)

	s.sse = sse.New()
	s.sse.AutoReplay = false
	s.sse.CreateStream("all")

	s.setupRoutes()
	go s.processingIncoming()

	// Start the HTTP server
	if err := http.ListenAndServe(s.opts.Addr, s.router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.log.Info("Stopping http server")
	return nil
}

func (s *Server) setupRoutes() {
	Static(s.router)
	if !s.opts.OfflineMode {
		s.router.HandleFunc("/messages", messageHandler(s.statusChannel, s.log))
	}
	s.router.Group(func(r chi.Router) {
		r.Get("/sse-events", s.sse.ServeHTTP)
	})

	s.router.HandleFunc("/clear-messages", s.clearMessages())

	s.router.HandleFunc("/", ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return HomePage(s.messages), nil
	}))
}

func (s *Server) processingIncoming() {
	s.log.Info("Processing incoming messages...")
	cnt := 1

	for {
		if s.opts.OfflineMode {
			cnt += 1
			val := fmt.Sprintf(`{"name":"event-%d"}`, cnt)
			msg := PubSubMessage{
				RawMessage: val,
				ReceivedAt: time.Now(),
			}
			s.log.Info("Received message", "message", msg)
			s.messages = append(s.messages, msg)
			s.broadcastNode(MessageRow(msg), "incoming-messages")

			channel := make(chan bool)
			// this is a goroutine which executes asynchronously
			go func() {
				time.Sleep(5 * time.Second)
				// send a message to the channel
				channel <- true
			}()

			// setup a channel listener
			select {
			case val := <-channel:
				s.log.Debug("Received value from channel", "val", val)
			}
		} else {
			for {
				msg := <-s.statusChannel
				s.log.Info("Received message", "message", msg)
				s.messages = append(s.messages, msg)
				s.broadcastNode(MessageRow(msg), "incoming-messages")
			}
		}
	}
}

func (s *Server) broadcastNode(data Node, event string) {
	var b bytes.Buffer
	if err := data.Render(&b); err != nil {
		s.log.Error("Failed to render node", "error", err)
		return
	}
	s.sse.Publish("all", &sse.Event{
		Event: []byte(event),
		Data:  b.Bytes(),
	})
}

func (s *Server) broadcastString(data string, event string) {
	s.sse.Publish("all", &sse.Event{
		Event: []byte(event),
		Data:  []byte(data),
	})
}
