package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	"github.com/go-chi/chi/v5"
)

var ErrConversationNotFound = errors.New("conversation not found")

// GetConversations godoc
//	@Summary		Get user conversations
//	@Description	Get all conversations for the authenticated user
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		store.Conversation
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/messages/conversations [get]
func (app *application) getConversationsHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)
	if user == nil {
		app.unauthorizedErrorResponse(w, r, ErrConversationNotFound)
		return
	}

	conversations, err := app.store.Messages.GetUserConversations(r.Context(), user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusOK, conversations); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetConversation godoc
//	@Summary		Get specific conversation
//	@Description	Get a conversation by ID
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			conversationID	path		int	true	"Conversation ID"
//	@Success		200				{object}	store.Conversation
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/messages/conversations/{conversationID} [get]
func (app *application) getConversationHandler(w http.ResponseWriter, r *http.Request) {
	convID, err := strconv.ParseInt(chi.URLParam(r, "conversationID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	conv, err := app.store.Messages.GetConversation(r.Context(), convID)
	if err != nil {
		if err == ErrConversationNotFound {
			app.notFoundResponse(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusOK, conv); err != nil {
		app.internalServerError(w, r, err)
	}
}

type CreateMessageRequest struct {
	Content string `json:"content"`
}


// CreateMessageHandler godoc
//
//	@Summary		Create message
//	@Description	Creates a new message in a conversation
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"Conversation ID"
//	@Param			req	body		CreateMessageRequest	true	"message"
//	@Success		201	{object}	store.Message			"Conversation created"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/messages/conversations/{id}/messages [post]
func (app *application) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from the context
	user := getUserfromCtx(r)
	if user == nil {
		WriteJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get the conversation ID from the URL
	idStr := chi.URLParam(r, "id")
	conversationID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid conversation ID"))
		return
	}

	// Parse the request body
	var req CreateMessageRequest
	err = ReadJSON(w, r, &req)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the message content
	if req.Content == "" {
		app.badRequestResponse(w, r, errors.New("content is required"))
		return
	}

	// Create the message
	message := &store.Message{
		ConversationID: conversationID,
		SenderID:       user.ID,
		Content:        req.Content,
	}

	err = app.store.Messages.CreateMessage(r.Context(), message)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("content is required"))
		return
	}

	WriteJSON(w, http.StatusCreated, message)
}

// GetMessages godoc
//	@Summary		Get conversation messages
//	@Description	Get messages in a conversation with pagination
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			conversationID	path		int	true	"Conversation ID"
//	@Param			limit			query		int	false	"Limit"		default(50)
//	@Param			offset			query		int	false	"Offset"	default(0)
//	@Success		200				{array}		store.Message
//	@Failure		400				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/messages/{conversationID} [get]
func (app *application) getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	convID, err := strconv.ParseInt(chi.URLParam(r, "conversationID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	messages, err := app.store.Messages.GetConversationMessages(r.Context(), convID, limit, offset)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusOK, messages); err != nil {
		app.internalServerError(w, r, err)
	}
}

// MarkConversationRead godoc
//	@Summary		Mark conversation as read
//	@Description	Mark all messages in a conversation as read for a user
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			conversationID	path		int		true	"Conversation ID"
//	@Success		204				{string}	string	"Conversation marked as read"
//	@Failure		400				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/messages/{conversationID}/read [put]
func (app *application) markConversationReadHandler(w http.ResponseWriter, r *http.Request) {
	convID, err := strconv.ParseInt(chi.URLParam(r, "conversationID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserfromCtx(r)
	if user == nil {
		app.unauthorizedErrorResponse(w,r,ErrConversationNotFound)
		return
	}

	err = app.store.Messages.MarkConversationAsRead(r.Context(), convID, user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetUnreadCount godoc
//	@Summary		Get unread message count
//	@Description	Get total number of unread messages for the user
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	object	"count":	int
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/messages/unread [get]
func (app *application) getUnreadCountHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)
	if user == nil {
		app.unauthorizedErrorResponse(w, r, ErrConversationNotFound)
		return
	}

	count, err := app.store.Messages.GetUnreadCount(r.Context(), user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusOK, map[string]int{"count": count}); err != nil {
		app.internalServerError(w, r, err)
	}
}

type CreateConversationRequest struct {
	OtherUserID int64 `json:"other_user_id"`
}

// CreateConversationHandler godoc
//
//	@Summary		Create conversation
//	@Description	Creates a new conversation with another user
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			req	body		CreateConversationRequest	true	"Create conversation request"
//	@Success		201	{object}	store.Conversation			"Conversation created"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/messages/conversations [post]
func (app *application) createConversationHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user ID from the context
	user := getUserfromCtx(r)
	if user == nil {
		WriteJSON(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse the request body
	var req CreateConversationRequest
	err := ReadJSON(w, r, &req)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the other user ID
	if req.OtherUserID == 0 {
		app.badRequestResponse(w, r, errors.New("other_user_id is required"))
		return
	}

	if req.OtherUserID == user.ID {
		app.badRequestResponse(w, r, errors.New("cannot create conversation with yourself"))
		return
	}

	// Verify the other user exists
	otherUser, err := app.store.Users.GetByID(r.Context(), req.OtherUserID)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			app.notFoundResponse(w, r,err)
			return
		}
		app.badRequestResponse(w, r, errors.New("cannot create conversation with yourself"))
		return
	}

	// Get or create the conversation
	conversation, err := app.store.Messages.GetOrCreateConversationByUsers(r.Context(), user.ID, otherUser.ID)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("cannot create conversation with yourself"))
		return
	}

	WriteJSON(w, http.StatusCreated, conversation)
}