package store

import (
	"context"
	"database/sql"
)

type SocialMedia struct {
	ID     int64  `json:"id"`
	Userid int64  `json:"userid"`
	Name   string `json:"name"`
	Link   string `json:"link"`
}

type SocialMediaStore struct {
	db *sql.DB
}

func (s *SocialMediaStore) CreateSocialMedia(ctx context.Context, socialmedia *SocialMedia) error {
	query := `INSERT INTO socialmedia (userid, name, link)
		VALUES ($1, $2, $3)
		RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.db.QueryRowContext(ctx, query, socialmedia.Userid, socialmedia.Name, socialmedia.Link).Scan(&socialmedia.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SocialMediaStore) GetSocialMediaById(ctx context.Context, id int64) (*SocialMedia, error) {
	query := `SELECT * FROM socialmedia WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	socialmedia := &SocialMedia{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&socialmedia.ID, &socialmedia.Userid, &socialmedia.Name, &socialmedia.Link,
	)
	if err != nil {
		return nil, err
	}

	return socialmedia, nil
}

func (s *SocialMediaStore) GetSocialMediaByUserId(ctx context.Context, userid int64) ([]*SocialMedia, error) {
	query := `SELECT * FROM socialmedia WHERE userid = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	socialmedias := []*SocialMedia{}
	for rows.Next() {
		socialmedia := &SocialMedia{}
		err := rows.Scan(
			&socialmedia.ID, &socialmedia.Userid, &socialmedia.Name, &socialmedia.Link,
		)
		if err != nil {
			return nil, err
		}
		socialmedias = append(socialmedias, socialmedia)
	}

	return socialmedias, nil
}

// update
func (s *SocialMediaStore) UpdateSocialMedia(ctx context.Context, socialmedia *SocialMedia) error {
	query := `UPDATE socialmedia SET name = $1, link = $2
	WHERE id = $3`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, socialmedia.Name, socialmedia.Link, socialmedia.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// delete
func (s *SocialMediaStore) DeleteSocialMedia(ctx context.Context, socialmediaID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		DELETE FROM socialmedia WHERE id = $1
	`, socialmediaID)
	return err
}
