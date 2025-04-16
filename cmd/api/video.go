package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/pion/webrtc/v3"
)

type PeerConnection struct {
	ID          string
	PC          *webrtc.PeerConnection
	DataChannel *webrtc.DataChannel
	RoomID      string
	IsInitiator bool
}

// Room represents a meeting room with participants
type Room struct {
	ID          string
	Connections map[string]*PeerConnection
	mutex       sync.Mutex
}

// SignalMessage represents a WebRTC signaling message
type SignalMessage struct {
	Type      string                   `json:"type"`
	SDP       string                   `json:"sdp,omitempty"`
	Candidate *webrtc.ICECandidateInit `json:"candidate,omitempty"`
	UserID    string                   `json:"userId"`
	RoomID    string                   `json:"roomId"`
	To        string                   `json:"to,omitempty"`
	From      string                   `json:"from,omitempty"` // Add this field
}

// Global state
var (
	rooms    = make(map[string]*Room)
	roomLock sync.Mutex
)

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	var roomRequest struct {
		RoomID string `json:"roomId"`
		UserID string `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&roomRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	roomLock.Lock()
	defer roomLock.Unlock()

	if _, exists := rooms[roomRequest.RoomID]; exists {
		http.Error(w, "Room already exists", http.StatusBadRequest)
		return
	}

	newRoom := &Room{
		ID:          roomRequest.RoomID,
		Connections: make(map[string]*PeerConnection),
	}
	rooms[roomRequest.RoomID] = newRoom

	// Create peer connection for the room creator
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		http.Error(w, "Failed to create peer connection", http.StatusInternalServerError)
		return
	}

	// Create connection object
	connection := &PeerConnection{
		ID:          roomRequest.UserID,
		PC:          peerConnection,
		RoomID:      roomRequest.RoomID,
		IsInitiator: true,
	}

	// Configure ICE handlers
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		candidateJSON, err := json.Marshal(candidate.ToJSON())
		if err != nil {
			log.Println("Failed to marshal ICE candidate:", err)
			return
		}

		log.Println("ICE candidate found:", string(candidateJSON))
		// We would usually send this to the client
	})

	// Add to room
	newRoom.Connections[roomRequest.UserID] = connection

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"roomId": roomRequest.RoomID,
		"userId": roomRequest.UserID,
		"status": "created",
	})
}

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")

	var joinRequest struct {
		UserID string `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&joinRequest); err != nil {
		http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
		return
	}

	roomLock.Lock()
	defer roomLock.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		http.Error(w, "Room not found: "+roomID, http.StatusNotFound)
		return
	}

	// Check if user is already in the room
	if _, userExists := room.Connections[joinRequest.UserID]; userExists {
		// User is already in the room, return success with the other peer ID
		var otherPeerID string
		for id := range room.Connections {
			if id != joinRequest.UserID {
				otherPeerID = id
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"roomId":      roomID,
			"userId":      joinRequest.UserID,
			"otherPeerId": otherPeerID,
			"status":      "rejoined",
		})
		return
	}

	if len(room.Connections) >= 2 {
		http.Error(w, "Room is full (maximum 2 participants allowed)", http.StatusForbidden)
		return
	}

	// Create peer connection for joining user
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		http.Error(w, "Failed to create peer connection", http.StatusInternalServerError)
		return
	}

	// Create connection object
	connection := &PeerConnection{
		ID:          joinRequest.UserID,
		PC:          peerConnection,
		RoomID:      roomID,
		IsInitiator: false,
	}

	// Configure ICE handlers
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		candidateJSON, err := json.Marshal(candidate.ToJSON())
		if err != nil {
			log.Println("Failed to marshal ICE candidate:", err)
			return
		}

		log.Println("ICE candidate found:", string(candidateJSON))
		// We would usually send this to the client
	})

	// Add to room
	room.Connections[joinRequest.UserID] = connection

	// Get other peer's ID
	var otherPeerID string
	for id := range room.Connections {
		if id != joinRequest.UserID {
			otherPeerID = id
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"roomId":      roomID,
		"userId":      joinRequest.UserID,
		"otherPeerId": otherPeerID,
		"status":      "joined",
	})
}

func handleOffer(room *Room, signal SignalMessage, w http.ResponseWriter) {
	// Find the connection for the recipient
	recipientConn, exists := room.Connections[signal.To]
	if !exists {
		http.Error(w, "Recipient not found in room", http.StatusNotFound)
		return
	}

	// Create a session description from the offer
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  signal.SDP,
	}

	// Set the remote description
	if err := recipientConn.PC.SetRemoteDescription(offer); err != nil {
		http.Error(w, "Failed to set remote description", http.StatusInternalServerError)
		return
	}

	// Create answer
	answer, err := recipientConn.PC.CreateAnswer(nil)
	if err != nil {
		http.Error(w, "Failed to create answer", http.StatusInternalServerError)
		return
	}

	// Set local description
	if err = recipientConn.PC.SetLocalDescription(answer); err != nil {
		http.Error(w, "Failed to set local description", http.StatusInternalServerError)
		return
	}

	// Return the answer to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"type": "answer",
		"sdp":  answer.SDP,
		"from": signal.To,
		"to":   signal.UserID,
	})
}

func handleAnswer(room *Room, signal SignalMessage, w http.ResponseWriter) {
	// Find the connection for the recipient
	recipientConn, exists := room.Connections[signal.To]
	if !exists {
		http.Error(w, "Recipient not found in room", http.StatusNotFound)
		return
	}

	// Create a session description from the answer
	answer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  signal.SDP,
	}

	// Set the remote description
	if err := recipientConn.PC.SetRemoteDescription(answer); err != nil {
		http.Error(w, "Failed to set remote description", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func handleCandidate(room *Room, signal SignalMessage, w http.ResponseWriter) {
	// Find the connection for the recipient
	recipientConn, exists := room.Connections[signal.To]
	if !exists {
		http.Error(w, "Recipient not found in room", http.StatusNotFound)
		return
	}

	// Add the ICE candidate
	if err := recipientConn.PC.AddICECandidate(*signal.Candidate); err != nil {
		http.Error(w, "Failed to add ICE candidate", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

var pendingSignals = make(map[string][]SignalMessage)
var pendingSignalsLock sync.Mutex

func roomStatusHandler(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")

	roomLock.Lock()
	defer roomLock.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Get all connection IDs in the room
	var connectionIDs []string
	for id := range room.Connections {
		connectionIDs = append(connectionIDs, id)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"roomId":      roomID,
		"connections": connectionIDs,
	})
}

func pendingSignalsHandler(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	userID := chi.URLParam(r, "userID")

	pendingSignalsLock.Lock()
	defer pendingSignalsLock.Unlock()

	// Get key for this room and user
	key := roomID + "_" + userID

	// Get all pending signals for this user
	signals := pendingSignals[key]

	// Clear the signals after retrieving them
	delete(pendingSignals, key)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signals)
}

// Modify the signal handler to store signals for polling
func signalHandler(w http.ResponseWriter, r *http.Request) {
	var signal SignalMessage
	if err := json.NewDecoder(r.Body).Decode(&signal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	roomLock.Lock()
	room, exists := rooms[signal.RoomID]
	roomLock.Unlock()

	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Add the "from" field to the signal
	signal.From = signal.UserID

	// Process the signal based on its type
	switch signal.Type {
	case "offer":
		// Store the offer for the recipient to poll
		pendingSignalsLock.Lock()
		key := signal.RoomID + "_" + signal.To
		pendingSignals[key] = append(pendingSignals[key], signal)
		pendingSignalsLock.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "offer_queued"})

	case "answer":
		// For answers, we can respond directly since the initiator is waiting
		handleAnswer(room, signal, w)

	case "candidate":
		// Store the ICE candidate for the recipient to poll
		pendingSignalsLock.Lock()
		key := signal.RoomID + "_" + signal.To
		pendingSignals[key] = append(pendingSignals[key], signal)
		pendingSignalsLock.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "candidate_queued"})

	default:
		http.Error(w, "Unknown signal type", http.StatusBadRequest)
	}
}
