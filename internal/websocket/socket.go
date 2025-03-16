package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/Althaf66/Appointr/internal/store"
	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	connections map[int64]map[*websocket.Conn]bool // conversationID -> connections
	mutex       sync.RWMutex
	store       store.Storage
	Upgrader    websocket.Upgrader
}

func NewWebSocketManager(store store.Storage) *WebSocketManager {
	return &WebSocketManager{
		connections: make(map[int64]map[*websocket.Conn]bool),
		store:       store,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// RegisterConnection adds a new WebSocket connection for a conversation
func (wm *WebSocketManager) RegisterConnection(conversationID int64, conn *websocket.Conn) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	if _, exists := wm.connections[conversationID]; !exists {
		wm.connections[conversationID] = make(map[*websocket.Conn]bool)
	}
	wm.connections[conversationID][conn] = true
}

// RemoveConnection removes a WebSocket connection
func (wm *WebSocketManager) RemoveConnection(conversationID int64, conn *websocket.Conn) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	if conns, exists := wm.connections[conversationID]; exists {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(wm.connections, conversationID)
		}
	}
}

// BroadcastMessage sends a message to all connections in a conversation
func (wm *WebSocketManager) BroadcastMessage(conversationID int64, message *store.Message) error {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	conns, exists := wm.connections[conversationID]
	if !exists {
		return nil
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for conn := range conns {
		err := conn.WriteMessage(websocket.TextMessage, messageJSON)
		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
			// Consider removing the connection if it's broken
			go wm.RemoveConnection(conversationID, conn)
		}
	}
	return nil
}
