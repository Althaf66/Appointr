package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

type Hub struct {
	// Registered clients by conversation ID
	clients map[int]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Messages to broadcast to specific conversation
	broadcast chan *Message

	// Mutex for concurrent access to clients map
	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.ConversationID]; !ok {
				h.clients[client.ConversationID] = make(map[*Client]bool)
			}
			h.clients[client.ConversationID][client] = true
			h.mu.Unlock()
			log.Printf("Client %d registered for conversation %d", client.ID, client.ConversationID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ConversationID]; ok {
				if _, exists := h.clients[client.ConversationID][client]; exists {
					delete(h.clients[client.ConversationID], client)
					close(client.Send)
					// If no more clients in this conversation, clean up
					if len(h.clients[client.ConversationID]) == 0 {
						delete(h.clients, client.ConversationID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client %d unregistered from conversation %d", client.ID, client.ConversationID)

		case message := <-h.broadcast:
			h.mu.Lock()
			if clients, ok := h.clients[message.ConversationID]; ok {
				// Marshal the message to JSON
				data, err := json.Marshal(message)
				if err != nil {
					log.Printf("Error marshaling message: %v", err)
					h.mu.Unlock()
					continue
				}

				// Send to all clients in this conversation
				for client := range clients {
					select {
					case client.Send <- data:
					default:
						// If client's send buffer is full, remove the client
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Broadcasted message to conversation %d", message.ConversationID)
		}
	}
}
