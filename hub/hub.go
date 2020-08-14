package hub

import (
	"encoding/json"

	"github.com/nicpottier/decent/parser"
)

// Hub maintains the set of active clients and broadcasts messages to each client
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan string

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func New() *Hub {
	return &Hub{
		Broadcast:  make(chan string),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
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
