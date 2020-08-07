package main

import (
	"bytes"
	"encoding/json"

	"github.com/nicpottier/decent/parser"
)

// Hub maintains the set of active clients and broadcasts messages to each client
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
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
			mj := ParseToken(message)

			for client := range h.clients {
				select {
				case client.send <- mj:
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

func ParseToken(t []byte) []byte {
	mt, mb, err := parser.ReadNextToken(bytes.NewReader(t))
	if err != nil {
		return errorJSON(err)
	}

	m, err := parser.ParseMessage(mt, mb)
	if err != nil {
		return errorJSON(err)
	}

	b, err := json.Marshal(m)
	if err != nil {
		return errorJSON(err)
	}

	return b

}
