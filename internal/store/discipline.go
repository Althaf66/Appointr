package store

import (
	"context"
	"database/sql"
	"errors"
)

type Discipline struct {
	ID       int64  `json:"id"`
	Field    string `json:"field"`
	SubField string `json:"subfield"`
}

var ErrDisciplineNotFound = errors.New("field was not found")

type DisciplineStore struct {
	db *sql.DB
}

func (s *DisciplineStore) Create(ctx context.Context, dis *Discipline) error {
	query := `INSERT INTO discipline (field,subfield) VALUES ($1,$2) 
	RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, dis.Field, dis.SubField).
		Scan(&dis.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *DisciplineStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM discipline WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrDisciplineNotFound
	}
	return nil
}

func (s *DisciplineStore) Update(ctx context.Context, dis *Discipline) error {
	query := `UPDATE discipline SET field = $1, subfield = $2 
	WHERE id = $3 RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, dis.Field, dis.SubField, dis.ID).Scan(&dis.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrDisciplineNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *DisciplineStore) Get(ctx context.Context) ([]*Discipline, error) {
	query := `SELECT id, field, subfield FROM discipline`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var disciplineList []*Discipline
	for rows.Next() {
		var dis Discipline
		err := rows.Scan(&dis.ID, &dis.Field, &dis.SubField)
		if err != nil {
			return nil, err
		}
		disciplineList = append(disciplineList, &dis)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return disciplineList, nil
}

func (s *DisciplineStore) GetByField(ctx context.Context, field string) (*Discipline, error) {
	query := `
		SELECT id, field, subfield FROM discipline WHERE field = $1
	`
	var dis Discipline
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, field).Scan(&dis.ID, &dis.Field, &dis.SubField)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrDisciplineNotFound
		default:
			return nil, err
		}
	}
	return &dis, nil
}

func (s *DisciplineStore) GetByID(ctx context.Context, id int64) (*Discipline, error) {
	query := `
		SELECT id, field, subfield FROM discipline WHERE id = $1
	`
	var dis Discipline
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, id).Scan(&dis.ID, &dis.Field, &dis.SubField)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrDisciplineNotFound
		default:
			return nil, err
		}
	}
	return &dis, nil
}

// complete discipline api and its interface
