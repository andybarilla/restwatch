package main

import (
	"fmt"
	"log/slog"
	"sync"
)

const SSE_ALL_CLIENTS = "all"

type SSEConn struct {
	mu      sync.Mutex
	logger  *slog.Logger
	clients map[string][]chan string
}

func NewSSEConn(l *slog.Logger) *SSEConn {
	return &SSEConn{
		logger:  l,
		clients: make(map[string][]chan string),
	}
}

func (p *SSEConn) addClient(id string) *chan string {
	p.mu.Lock()
	defer func() {
		p.logger.Debug(fmt.Sprintf("Clients in add: %v", p.clients))
		for k, v := range p.clients {
			p.logger.Debug(fmt.Sprintf("Key: %s, value: %d", k, len(v)))
			p.logger.Debug(fmt.Sprintf("Channels from id=%s: %v", id, v))
		}
		p.mu.Unlock()
	}()

	c, ok := p.clients[id]
	if !ok {
		client := []chan string{make(chan string)}
		p.clients[id] = client
		return &client[0]
	}

	newCh := make(chan string)
	p.clients[id] = append(c, newCh)
	return &newCh
}

func (p *SSEConn) removeClient(id string, conn chan string) {
	p.mu.Lock()
	defer func() {
		p.logger.Debug(fmt.Sprintf("Clients in remove: %v", p.clients))
		for k, v := range p.clients {
			p.logger.Debug(fmt.Sprintf("Key: %s, value: %d", k, len(v)))
		}
		p.mu.Unlock()
	}()

	c, ok := p.clients[id]
	if !ok {
		return
	}

	pos := -1

	for i, ch := range c {
		if ch == conn {
			pos = i
		}
	}

	if pos == -1 {
		return
	}

	close(c[pos])
	c = append(c[:pos], c[pos+1:]...)
	if pos == 0 {
		delete(p.clients, id)
	}
}

func (p *SSEConn) broadcast(id string, data, event string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	c, ok := p.clients[id]
	if !ok {
		p.logger.Info(fmt.Sprintf("No clients for id=%s", id))
		return
	}

	p.logger.Info(fmt.Sprintf("Sending message to %d client(s)", len(c)))
	for _, ch := range c {
		ch <- fmt.Sprintf("event: %s\ndata: %s", event, data)
	}
}
