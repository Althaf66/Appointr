package store

import (
	"context"
	"database/sql"
)

type Education struct {
	ID        int    `json:"id"`
	Userid    int64  `json:"userid"`
	Year_from string `json:"year_from"`
	Year_to   string `json:"year_to"`
	Degree    string `json:"degree"`
	Field     string `json:"field"`
	Institute string `json:"institute"`
}

type EducationStore struct {
	db *sql.DB
}

func (s *EducationStore) CreateEducation(ctx context.Context, education *Education) error {
	query := `INSERT INTO education (userid, year_from, year_to, degree, field, institute)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.db.QueryRowContext(ctx, query, education.Userid, education.Year_from, education.Year_to,
		education.Degree, education.Field, education.Institute).Scan(&education.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *EducationStore) GetEducationById(ctx context.Context, id int64) (*Education, error) {
	query := `SELECT * FROM education WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	education := &Education{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&education.ID, &education.Userid, &education.Year_from, &education.Year_to, &education.Degree, &education.Field, &education.Institute,
	)
	if err != nil {
		return nil, err
	}

	return education, nil
}

func (s *EducationStore) GetEducationByUserId(ctx context.Context, userid int64) ([]*Education, error) {
	query := `SELECT * FROM education WHERE userid = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	educations := []*Education{}
	for rows.Next() {
		education := &Education{}
		err := rows.Scan(
			&education.ID, &education.Userid, &education.Year_from, &education.Year_to, &education.Degree, &education.Field, &education.Institute,
		)
		if err != nil {
			return nil, err
		}
		educations = append(educations, education)
	}

	return educations, nil
}

// update
func (s *EducationStore) UpdateEducation(ctx context.Context, education *Education) error {
	query := `UPDATE education SET year_from = $1, year_to = $2, degree = $3, field = $4, institute = $5
	WHERE id = $6`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, education.Year_from, education.Year_to, education.Degree,
		education.Field, education.Institute, education.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// delete
func (s *EducationStore) DeleteEducation(ctx context.Context, educationID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		DELETE FROM education WHERE id = $1
	`, educationID)
	return err
}
