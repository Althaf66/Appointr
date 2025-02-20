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
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UserStore{db},
		Expertise: &ExpertiseStore{db},
		Discipline: &DisciplineStore{db},
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
