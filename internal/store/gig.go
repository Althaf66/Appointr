package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Gig struct {
	ID          int64    `json:"id"`
	Userid      int64    `json:"userid"`
	Title       string   `json:"title"`
	Amount      float64  `json:"amount"`
	Description string   `json:"description"`
	Expertise   string   `json:"expertise"`
	Discipline  []string `json:"discipline"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type GigStore struct {
	db *sql.DB
}

func (s *GigStore) CreateGig(ctx context.Context, gig *Gig) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
		INSERT INTO gigs (userid, title, amount, description, expertise, discipline)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`,
		gig.Userid, gig.Title, gig.Amount, gig.Description, gig.Expertise, pq.Array(gig.Discipline)).
		Scan(&gig.ID, &gig.CreatedAt, &gig.UpdatedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *GigStore) GetAllGigs(ctx context.Context, limit, offset int) ([]*Gig, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	if limit <= 0 {
		limit = 50 // Default limit
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, userid, title, amount,description, expertise, discipline FROM gigs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gigs := []*Gig{}
	for rows.Next() {
		gig := &Gig{}
		err := rows.Scan(
			&gig.ID, &gig.Userid, &gig.Title, &gig.Amount, &gig.Description, &gig.Expertise, pq.Array(&gig.Discipline),
		)
		if err != nil {
			return nil, err
		}
		gigs = append(gigs, gig)
	}

	return gigs, nil
}

func (s *GigStore) GetGigsByExpertise(ctx context.Context, expertise string) ([]*Gig, error) {
	query := `SELECT id,userid,title,amount,expertise,discipline,created_at,updated_at FROM 
	gigs WHERE expertise ILIKE $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, "%"+expertise+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gigs := []*Gig{}
	for rows.Next() {
		gig := &Gig{}
		err := rows.Scan(
			&gig.ID, &gig.Userid, &gig.Title, &gig.Amount, &gig.Expertise, pq.Array(&gig.Discipline),
			&gig.CreatedAt, &gig.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		gigs = append(gigs, gig)
	}

	return gigs, nil
}

func (s *GigStore) GetGigByID(ctx context.Context, id int64) (*Gig, error) {
	query := `SELECT * FROM gigs WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var gig Gig
	err := s.db.QueryRowContext(ctx, query, id).Scan(&gig.ID, &gig.Userid, &gig.Title, &gig.Amount, &gig.Description, &gig.Expertise, pq.Array(&gig.Discipline),
		&gig.CreatedAt, &gig.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &gig, nil
}

// getbydiscipline
// getbymentorid
// getbyuserid

func (s *GigStore) UpdateGig(ctx context.Context, gig *Gig) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
        UPDATE gigs 
        SET title = $1, amount=$2, description = $3, expertise = $4, discipline = $5, updated_at = NOW()
        WHERE id = $6
    `, gig.Title, gig.Amount, gig.Description, gig.Expertise, pq.Array(gig.Discipline), gig.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *GigStore) DeleteGig(ctx context.Context, gigID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		DELETE FROM gigs WHERE id = $1
	`, gigID)
	return err
}
