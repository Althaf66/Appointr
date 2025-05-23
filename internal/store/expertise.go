package store

import (
	"context"
	"database/sql"
	"errors"
)

type Expertise struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Icon_svg string `json:"icon_svg"`
}

type Country struct {
	ID       int64  `json:"id"`
	CountryName     string `json:"countryname"`
}

var ErrExpertiseNotFound = errors.New("field was not found")

type ExpertiseStore struct {
	db *sql.DB
}

type CountryStore struct {
	db *sql.DB
}

func (s *ExpertiseStore) Create(ctx context.Context, exp *Expertise) error {
	query := `INSERT INTO expertise (name,icon_svg) VALUES ($1,$2) 
	RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, exp.Name, exp.Icon_svg).Scan(&exp.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ExpertiseStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM expertise WHERE id = $1`

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
		return ErrExpertiseNotFound
	}
	return nil
}

func (s *ExpertiseStore) Update(ctx context.Context, exp *Expertise) error {
	query := `UPDATE expertise SET name = $1, icon_svg = $2 
	WHERE id = $3 RETURNING id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, exp.Name, exp.Icon_svg, exp.ID).Scan(&exp.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrExpertiseNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *ExpertiseStore) Get(ctx context.Context) ([]*Expertise, error) {
	query := `SELECT id, name, icon_svg FROM expertise`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expertiseList []*Expertise
	for rows.Next() {
		var exp Expertise
		err := rows.Scan(&exp.ID, &exp.Name, &exp.Icon_svg)
		if err != nil {
			return nil, err
		}
		expertiseList = append(expertiseList, &exp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return expertiseList, nil
}

func (s *ExpertiseStore) GetByID(ctx context.Context, id int64) (*Expertise, error) {
	query := `
		SELECT id, name, icon_svg FROM expertise WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var exp Expertise
	err := s.db.QueryRowContext(ctx, query, id).Scan(&exp.ID, &exp.Name, &exp.Icon_svg)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrExpertiseNotFound
		default:
			return nil, err
		}
	}
	return &exp, nil
}

func (s *CountryStore) GetCountry(ctx context.Context) ([]*Country, error) {
	query := `SELECT countryname FROM countries`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var CountryList []*Country
	for rows.Next() {
		var country Country
		err := rows.Scan(&country.CountryName)
		if err != nil {
			return nil, err
		}
		CountryList = append(CountryList, &country)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return CountryList, nil
}
