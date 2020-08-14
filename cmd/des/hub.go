package main

import (
	"encoding/json"

	"github.com/nicpottier/decent/parser"
)

// Hub maintains the set of active clients and broadcasts messages to each client
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan string

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			m, err := parser.ParseMessage(message)
			if err != nil {
				m = errorJSON(err)
			}

			b, err := json.Marshal(m)
			if err != nil {
				b = errorJSON(err)
			}

			for client := range h.clients {
				select {
				case client.send <- b:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func errorJSON(err error) []byte {
	b, _ := json.Marshal(map[string]string{
		"type":  "error",
		"error": err.Error(),
	})
	return b
}
