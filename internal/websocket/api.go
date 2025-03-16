package websocket

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
)

func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conversationID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("upgrade successfully")
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// Get user ID from query param (simplified, you would use proper JWT validation in production)
	userID := 1
	userIDStr := r.URL.Query().Get("userId")
	if userIDStr != "" {
		userID, _ = strconv.Atoi(userIDStr)
	}

	// Upgrade HTTP connection to WebSocket

	// Create new client
	client := &Client{
		ID:             userID,
		Conn:           conn,
		ConversationID: conversationID,
		Send:           make(chan []byte, 256),
	}

	// Register client with hub
	h.register <- client

	// Start goroutines for reading and writing
	go client.writePump(h)
	go client.readPump(h)
}
