package store

import (
	"context"
	"database/sql"
)

type Experience struct {
	ID          int64  `json:"id"`
	Userid      int64  `json:"userid"`
	Year_from   string `json:"year_from"`
	Year_to     string `json:"year_to"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Description string `json:"description"`
}

type ExperienceStore struct {
	db *sql.DB
}

func (s *ExperienceStore) CreateExperience(ctx context.Context, experience *Experience) error {
	query := `INSERT INTO experience (userid, year_from, year_to, title, company, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.db.QueryRowContext(ctx, query, experience.Userid, experience.Year_from, experience.Year_to,
		experience.Title, experience.Company, experience.Description).Scan(&experience.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ExperienceStore) GetExperienceById(ctx context.Context, id int64) (*Experience, error) {
	query := `SELECT * FROM experience WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	experience := &Experience{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&experience.ID, &experience.Userid, &experience.Year_from, &experience.Year_to, &experience.Title, &experience.Company, &experience.Description,
	)
	if err != nil {
		return nil, err
	}

	return experience, nil
}

func (s *ExperienceStore) GetExperienceByUserId(ctx context.Context, userid int64) ([]*Experience, error) {
	query := `SELECT * FROM experience WHERE userid = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	experiences := []*Experience{}
	for rows.Next() {
		experience := &Experience{}
		err := rows.Scan(
			&experience.ID, &experience.Userid, &experience.Year_from, &experience.Year_to, &experience.Title, &experience.Company, &experience.Description,
		)
		if err != nil {
			return nil, err
		}
		experiences = append(experiences, experience)
	}

	return experiences, nil
}

// update
func (s *ExperienceStore) UpdateExperience(ctx context.Context, experience *Experience) error {
	query := `UPDATE experience SET year_from = $1, year_to = $2, title = $3, company = $4, description = $5
	WHERE id = $6`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, experience.Year_from, experience.Year_to, experience.Title,
		experience.Company, experience.Description, experience.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// delete
func (s *ExperienceStore) DeleteExperience(ctx context.Context, experienceID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		DELETE FROM experience WHERE id = $1
	`, experienceID)
	return err
}
