package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID             int
	Conn           *websocket.Conn
	ConversationID int
	Send           chan []byte
}

type Message struct {
	ID             int       `json:"id"`
	Content        string    `json:"content"`
	SenderID       int       `json:"senderId"`
	ConversationID int       `json:"conversationId"`
	CreatedAt      time.Time `json:"created_at"`
	Sender         User      `json:"sender"`
}

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512 * 1024) // Max message size
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Parse message from client
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Add sender details and timestamp
		msg.CreatedAt = time.Now()
		msg.Sender = User{
			ID:       c.ID,
			Username: fmt.Sprintf("User %d", c.ID), // In a real app, fetch from DB
		}

		// Save to database (simulated)
		msg.ID = int(time.Now().UnixNano() % 10000) // Simulate ID generation

		// Broadcast to all clients in the conversation
		h.broadcast <- &msg
	}
}

// Send messages to the client
func (c *Client) writePump(h *Hub) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
