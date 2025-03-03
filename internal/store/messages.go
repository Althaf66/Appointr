package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Conversation struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Derived fields (not stored directly in conversations table)
	Participants []*User     `json:"participants,omitempty"`
	LastMessage  *Message    `json:"last_message,omitempty"`
	OtherUser    *User       `json:"other_user,omitempty"` // The user that is not the current user
	Unread       int         `json:"unread"`              // Number of unread messages
}

// ConversationParticipant joins users to conversations
type ConversationParticipant struct {
	ConversationID int64 `json:"conversation_id"`
	UserID         int64 `json:"user_id"`
}

// Message represents a single message within a conversation
type Message struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversation_id"`
	SenderID       int64     `json:"sender_id"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
	IsRead         bool      `json:"is_read"`
	// Derived fields
	Sender *User `json:"sender,omitempty"`
}

// MessageStore implements MessageStorage interface
type MessageStore struct {
	db *sql.DB
}

// CreateConversation creates a new conversation between two users
func (s *MessageStore) CreateConversation(ctx context.Context, userID1, userID2 int64) (*Conversation, error) {
	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create the conversation
	var conversationID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO conversations DEFAULT VALUES
		RETURNING id, created_at, updated_at
	`).Scan(&conversationID, &time.Time{}, &time.Time{})
	if err != nil {
		return nil, err
	}

	// Add participants
	_, err = tx.ExecContext(ctx, `
		INSERT INTO conversation_participants (conversation_id, user_id)
		VALUES ($1, $2), ($1, $3)
	`, conversationID, userID1, userID2)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Return the new conversation
	return s.GetConversation(ctx, conversationID)
}

// GetConversation retrieves a conversation by ID with participants
func (s *MessageStore) GetConversation(ctx context.Context, conversationID int64) (*Conversation, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get the conversation
	var conv Conversation
	err := s.db.QueryRowContext(ctx, `
		SELECT id, created_at, updated_at
		FROM conversations
		WHERE id = $1
	`, conversationID).Scan(&conv.ID, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("conversation not found")
		}
		return nil, err
	}

	// Get participants
	rows, err := s.db.QueryContext(ctx, `
		SELECT u.id, u.username, u.email, u.created_at
		FROM users u
		JOIN conversation_participants cp ON u.id = cp.user_id
		WHERE cp.conversation_id = $1 AND u.is_active = true
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conv.Participants = []*User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		conv.Participants = append(conv.Participants, &user)
	}

	// Get last message
	// err = s.db.QueryRowContext(ctx, `
	// 	SELECT id, sender_id, content, created_at, is_read
	// 	FROM messages
	// 	WHERE conversation_id = $1
	// 	ORDER BY created_at DESC
	// 	LIMIT 1
	// `, conversationID).Scan(
	// 	&conv.LastMessage.ID,
	// 	&conv.LastMessage.SenderID,
	// 	&conv.LastMessage.Content,
	// 	&conv.LastMessage.CreatedAt,
	// 	&conv.LastMessage.IsRead,
	// )
	// if err != nil && err != sql.ErrNoRows {
	// 	return nil, err
	// }
	// It's okay if there's no last message (new conversation)

	return &conv, nil
}

// GetOrCreateConversationByUsers finds or creates a conversation between two users
func (s *MessageStore) GetOrCreateConversationByUsers(ctx context.Context, userID1, userID2 int64) (*Conversation, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Try to find existing conversation between these users
	var conversationID int64
	err := s.db.QueryRowContext(ctx, `
		SELECT cp1.conversation_id
		FROM conversation_participants cp1
		JOIN conversation_participants cp2 ON cp1.conversation_id = cp2.conversation_id
		WHERE cp1.user_id = $1 AND cp2.user_id = $2
	`, userID1, userID2).Scan(&conversationID)

	if err != nil {
		if err == sql.ErrNoRows {
			// No conversation exists, create a new one
			return s.CreateConversation(ctx, userID1, userID2)
		}
		return nil, err
	}

	// Conversation exists, return it
	return s.GetConversation(ctx, conversationID)
}

// GetUserConversations retrieves all conversations for a user
func (s *MessageStore) GetUserConversations(ctx context.Context, userID int64) ([]*Conversation, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get all conversation IDs for this user
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.created_at, c.updated_at,
			(SELECT COUNT(*) FROM messages m WHERE m.conversation_id = c.id AND m.sender_id != $1 AND m.is_read = false) as unread_count,
			m.id, m.sender_id, m.content, m.created_at, m.is_read,
			u.id, u.username, u.email, u.created_at
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		JOIN conversation_participants cp2 ON c.id = cp2.conversation_id AND cp2.user_id != $1
		JOIN users u ON cp2.user_id = u.id
		LEFT JOIN messages m ON m.id = (
			SELECT id FROM messages 
			WHERE conversation_id = c.id 
			ORDER BY created_at DESC 
			LIMIT 1
		)
		WHERE cp.user_id = $1 AND u.is_active = true
		ORDER BY c.updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversations := []*Conversation{}
	for rows.Next() {
		conv := &Conversation{
			LastMessage: &Message{},
			OtherUser:   &User{},
		}
		
		err := rows.Scan(
			&conv.ID, &conv.CreatedAt, &conv.UpdatedAt, &conv.Unread,
			&conv.LastMessage.ID, &conv.LastMessage.SenderID, &conv.LastMessage.Content, &conv.LastMessage.CreatedAt, &conv.LastMessage.IsRead,
			&conv.OtherUser.ID, &conv.OtherUser.Username, &conv.OtherUser.Email, &conv.OtherUser.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Handle nullable message fields
		// if msgID.Valid {
		// 	conv.LastMessage.ID = msgID.Int64
		// 	conv.LastMessage.SenderID = senderID.Int64
		// 	conv.LastMessage.Content = content.String
		// 	// Handle other fields
		// } else {
		// 	conv.LastMessage = nil
		// }
		
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// CreateMessage adds a new message to a conversation
func (s *MessageStore) CreateMessage(ctx context.Context, message *Message) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Verify the sender is a participant in the conversation
	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM conversation_participants
		WHERE conversation_id = $1 AND user_id = $2
	`, message.ConversationID, message.SenderID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("sender is not a participant in this conversation")
	}

	// Insert the message
	return s.db.QueryRowContext(ctx, `
		INSERT INTO messages (conversation_id, sender_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, message.ConversationID, message.SenderID, message.Content).Scan(&message.ID, &message.CreatedAt)
}

// GetConversationMessages retrieves messages for a conversation with pagination
func (s *MessageStore) GetConversationMessages(ctx context.Context, conversationID int64, limit, offset int) ([]*Message, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 50 // Default limit
	}

	// Query messages with sender information
	rows, err := s.db.QueryContext(ctx, `
		SELECT m.id, m.sender_id, m.content, m.created_at, m.is_read,
			   u.id, u.username, u.email, u.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.conversation_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*Message{}
	for rows.Next() {
		msg := &Message{
			ConversationID: conversationID,
			Sender:         &User{},
		}
		
		err := rows.Scan(
			&msg.ID, &msg.SenderID, &msg.Content, &msg.CreatedAt, &msg.IsRead,
			&msg.Sender.ID, &msg.Sender.Username, &msg.Sender.Email, &msg.Sender.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		messages = append(messages, msg)
	}

	return messages, nil
}

// MarkConversationAsRead marks all messages in a conversation as read for a user
func (s *MessageStore) MarkConversationAsRead(ctx context.Context, conversationID, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		UPDATE messages
		SET is_read = true
		WHERE conversation_id = $1 AND sender_id != $2 AND is_read = false
	`, conversationID, userID)
	
	return err
}

// GetUnreadCount gets the total number of unread messages for a user
func (s *MessageStore) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM messages m
		JOIN conversation_participants cp ON m.conversation_id = cp.conversation_id
		WHERE cp.user_id = $1 AND m.sender_id != $1 AND m.is_read = false
	`, userID).Scan(&count)
	
	return count, err
}