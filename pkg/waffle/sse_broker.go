package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// SSEBroker Maintains list of connected clients and channels to
// handle the addition and removal of clients
type SSEBroker struct {
	Notifier       chan []byte
	newClients     chan chan []byte
	closingClients chan chan []byte
	clients        map[chan []byte]bool
}

// Creates a new Handler and starts goroutine to handle incoming events and
// push them to connected clients
func startSSEBroker() *SSEBroker {
	broker := &SSEBroker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	go broker.listen()
	go broker.keepalive()

	return broker
}

// SSE endpoint handler
func (broker *SSEBroker) connHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a new channel for the connected client and add it to the list
	messageChan := make(chan []byte)
	broker.newClients <- messageChan

	// Deregister client when the surrounding function returns
	defer func() {
		broker.closingClients <- messageChan
	}()

	ctx := r.Context()

	// Done() returns a channel that is closed when work related to this
	// request context should be canceled. Spawn a goroutine that blocks while the
	// channel is open, and removes the client from the client list if and when
	// it is closed
	go func() {
		<-ctx.Done()
		broker.closingClients <- messageChan
	}()

	// Runs for the duration of the client connection
	for {
		// Write an SSE compatible message to the response writer
		// https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#Examples
		fmt.Fprintf(w, "data: %s\n\n", <-messageChan)

		// Immediately flush data to the client, don't buffer
		flusher.Flush()
	}
}

// Blocking function that takes events off a channel and handles them
func (broker *SSEBroker) listen() {
	for {
		select {
		case s := <-broker.newClients:
			log.Printf("[waffle/sse] Client connected: %v", s)
			broker.clients[s] = true
		case s := <-broker.closingClients:
			log.Printf("[waffle/sse] Client disconnected: %v", s)
			delete(broker.clients, s)
		case event := <-broker.Notifier:
			log.Printf(
				"[waffle/sse] New event. Publishing to clients: %v",
				broker.clients,
			)

			for client := range broker.clients {
				client <- event
			}
		}
	}
}

// Prevent request timeouts on Heroku
func (broker *SSEBroker) keepalive() {
	for {
		time.Sleep(time.Second * 10)
		broker.Notifier <- []byte("{\"keepalive\":true}")
	}
}
