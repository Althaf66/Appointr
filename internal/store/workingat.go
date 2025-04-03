package store

import (
	"context"
	"database/sql"
)

type WorkingAt struct {
	ID        int64  `json:"id"`
	Userid    int64  `json:"userid"`
	Title     string `json:"title"`
	Company   string `json:"company"`
	TotalYear int64  `json:"totalyear"`
	Month     int64  `json:"month"`
	Linkedin  string `json:"linkedin"`
	Github    string `json:"github"`
	Instagram string `json:"instagram"`
}

type WorkingAtStore struct {
	db *sql.DB
}

func (s *WorkingAtStore) CreateWorkingAt(ctx context.Context, workingat *WorkingAt) error {
	query := `INSERT INTO workingat (userid, title, company, totalyear, month, linkedin, github, instagram) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.db.QueryRowContext(ctx, query, workingat.Userid, workingat.Title, workingat.Company,
		workingat.TotalYear, workingat.Month, workingat.Linkedin, workingat.Github, workingat.Instagram).Scan(&workingat.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *WorkingAtStore) GetWorkingAtById(ctx context.Context, id int64) (*WorkingAt, error) {
	query := `SELECT * FROM workingat WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	workingat := &WorkingAt{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&workingat.ID, &workingat.Userid, &workingat.Title, &workingat.Company, &workingat.TotalYear, &workingat.Month,
		&workingat.Linkedin, &workingat.Github, &workingat.Instagram)
	if err != nil {
		return nil, err
	}

	return workingat, nil
}

func (s *WorkingAtStore) GetWorkingAtByUserId(ctx context.Context, userid int64) ([]*WorkingAt, error) {
	query := `SELECT * FROM workingat WHERE userid = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workingats := []*WorkingAt{}
	for rows.Next() {
		workingat := &WorkingAt{}
		err := rows.Scan(
			&workingat.ID, &workingat.Userid, &workingat.Title, &workingat.Company, &workingat.TotalYear, &workingat.Month,
			&workingat.Linkedin, &workingat.Github, &workingat.Instagram)
		if err != nil {
			return nil, err
		}
		workingats = append(workingats, workingat)
	}

	return workingats, nil
}

// update
func (s *WorkingAtStore) UpdateWorkingAt(ctx context.Context, workingat *WorkingAt) error {
	query := `UPDATE workingat SET title = $1, company = $2,  totalyear = $3, month = $4,
	linkedin = $5, github = $6, instagram = $7
	WHERE id = $5`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, workingat.Title, workingat.Company,
		workingat.TotalYear, workingat.Month, workingat.Linkedin, workingat.Github, workingat.Instagram, workingat.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// delete
func (s *WorkingAtStore) DeleteWorkingAt(ctx context.Context, workingatID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		DELETE FROM workingat WHERE id = $1
	`, workingatID)
	return err
}
