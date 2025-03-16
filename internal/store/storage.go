package store

import (
	"context"
	"database/sql"
	"time"
)

type Storage struct {
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		// GauthCreate(context.Context, string, string) (*User, error)
		GetByID(context.Context, int64) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
		GetByEmail(context.Context, string) (*User, error)
	}
	Expertise interface {
		Create(context.Context, *Expertise) error
		Get(context.Context) ([]*Expertise, error)
		Update(context.Context, *Expertise) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*Expertise, error)
	}
	Discipline interface {
		Create(context.Context, *Discipline) error
		Get(context.Context) ([]*Discipline, error)
		Update(context.Context, *Discipline) error
		Delete(context.Context, int64) error
		GetByID(context.Context, int64) (*Discipline, error)
		GetByField(context.Context, string) (*Discipline, error)
	}
	Messages interface {
		CreateConversation(ctx context.Context, userID1, userID2 int64) (*Conversation, error)
		GetConversation(ctx context.Context, conversationID int64) (*Conversation, error)
		GetOrCreateConversationByUsers(ctx context.Context, userID1, userID2 int64) (*Conversation, error)
		GetUserConversations(ctx context.Context, userID int64) ([]*Conversation, error)
		CreateMessage(ctx context.Context, message *Message) error
		GetConversationMessages(ctx context.Context, conversationID int64, limit, offset int) ([]*Message, error)
		MarkConversationAsRead(ctx context.Context, conversationID, userID int64) error
		GetUnreadCount(ctx context.Context, userID int64) (int, error)
	}
	Mentor interface {
		CreateMentor(ctx context.Context, mentor *Mentor) error
		GetAllMentors(ctx context.Context, limit, offset int) ([]*Mentor, error)
		GetMentorByName(ctx context.Context, name string) ([]*Mentor, error)
		GetMentorByID(ctx context.Context, id int64) (*Mentor, error)
		UpdateMentor(ctx context.Context, mentor *Mentor) error
		DeleteMentor(ctx context.Context, mentorID int64) error
		// GetMentorsByExpertise(ctx context.Context, expertise string) ([]*Mentor, error)
		// GetMentorsByDiscipline(ctx context.Context, discipline string) ([]*Mentor, error)
	}
	Gig interface {
		CreateGig(ctx context.Context, gig *Gig) error
		GetAllGigs(ctx context.Context, limit, offset int) ([]*Gig, error)
		GetGigsByExpertise(ctx context.Context, expertise string) ([]*Gig, error)
		UpdateGig(ctx context.Context, gig *Gig) error
		DeleteGig(ctx context.Context, gigID int64) error
		GetGigByID(ctx context.Context, id int64) (*Gig, error)
		// GetGigByMentorID(ctx context.Context, mentorID int64) ([]*Gig, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Users:      &UserStore{db},
		Expertise:  &ExpertiseStore{db},
		Discipline: &DisciplineStore{db},
		Messages:   &MessageStore{db},
		Mentor:     &MentorStore{db},
		Gig:        &GigStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
