package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	"github.com/Althaf66/Appointr/internal/websocket"
	chi "github.com/go-chi/chi/v5"
	ws "github.com/gorilla/websocket"
)

func (app *application) HandleWebSocket(wm *websocket.WebSocketManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get conversation ID from URL
		conversationID, err := strconv.ParseInt(chi.URLParam(r, "conversationID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := wm.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		defer conn.Close()

		// Register the connection
		wm.RegisterConnection(conversationID, conn)
		defer wm.RemoveConnection(conversationID, conn)

		// Handle incoming messages
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			// Parse incoming message
			var incomingMsg struct {
				Content  string `json:"content"`
				SenderID int64  `json:"senderId"`
			}
			if err := json.Unmarshal(msgBytes, &incomingMsg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Create message in database
			message := &store.Message{
				ConversationID: conversationID,
				SenderID:       incomingMsg.SenderID,
				Content:        incomingMsg.Content,
			}

			if err := app.store.Messages.CreateMessage(context.Background(), message); err != nil {
				log.Printf("Error creating message: %v", err)
				continue
			}

			// Get sender information
			sender, err := app.store.Users.GetByID(context.Background(), message.SenderID)
			if err != nil {
				log.Printf("Error getting sender: %v", err)
				continue
			}
			message.Sender = sender

			// Broadcast to all connected clients
			if err := wm.BroadcastMessage(conversationID, message); err != nil {
				log.Printf("Error broadcasting message: %v", err)
			}
		}
	}
}
