package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrNotFound = errors.New("resource not found")

type Mentor struct {
	ID               int64              `json:"id"`
	Userid           int64              `json:"userid"`
	Name             string             `json:"name"`
	Country          string             `json:"country"`
	Language         []string           `json:"language"`
	Gigs             []*Gig             `json:"gigs"`
	Education        []*Education       `json:"education"`
	Experience       []*Experience      `json:"experience"`
	WorkingNow       *WorkingAt         `json:"working_at"`
	SocialMedia      []*SocialMedia     `json:"social_media"`
	YearOfExperience *YearsOfExperience `json:"years_of_experience"`
	CreatedAt        string             `json:"created_at"`
	UpdatedAt        string             `json:"updated_at"`
}

type YearsOfExperience struct {
	Userid int64 `json:"userid"`
	Year   int64 `json:"year"`
	Month  int64 `json:"month"`
}

type SocialMedia struct {
	Userid int64  `json:"userid"`
	Name   string `json:"name"`
	Link   string `json:"link"`
}

type WorkingAt struct {
	Userid  int64  `json:"userid"`
	Title   string `json:"title"`
	Company string `json:"company"`
}

type Education struct {
	Userid    int64  `json:"userid"`
	Year_from string `json:"year_from"`
	Year_to   string `json:"year_to"`
	Degree    string `json:"degree"`
	Field     string `json:"field"`
	Institute string `json:"institute"`
}

type Experience struct {
	Userid      int64  `json:"userid"`
	Year_from   string `json:"year_from"`
	Year_to     string `json:"year_to"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Description string `json:"description"`
}

type MentorStore struct {
	db *sql.DB
}

// CreateMentor creates a new mentor
func (s *MentorStore) CreateMentor(ctx context.Context, mentor *Mentor) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert mentor basic details
	err = tx.QueryRowContext(ctx, `
		INSERT INTO mentors (userid, name, country, language)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`,
		mentor.Userid, mentor.Name, mentor.Country, pq.Array(mentor.Language)).Scan(
		&mentor.ID, &mentor.CreatedAt, &mentor.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetAllMentors retrieves all mentors
func (s *MentorStore) GetAllMentors(ctx context.Context, limit, offset int) ([]*Mentor, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	if limit <= 0 {
		limit = 50 // Default limit
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, userid, name, country, language, created_at, updated_at
		FROM mentors
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mentors := []*Mentor{}
	for rows.Next() {
		mentor := &Mentor{}
		err := rows.Scan(
			&mentor.ID, &mentor.Userid, &mentor.Name, &mentor.Country, pq.Array(&mentor.Language), &mentor.CreatedAt, &mentor.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		mentors = append(mentors, mentor)
	}

	return mentors, nil
}

// GetMentorByName finds mentors by name (partial match)
func (s *MentorStore) GetMentorByName(ctx context.Context, name string) ([]*Mentor, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, userid, name, country, language, created_at, updated_at
		FROM mentors
		WHERE name ILIKE $1
		ORDER BY name`, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mentors := []*Mentor{}
	for rows.Next() {
		mentor := &Mentor{}
		err := rows.Scan(
			&mentor.ID, &mentor.Userid, &mentor.Name, &mentor.Country,
			pq.Array(&mentor.Language), &mentor.CreatedAt, &mentor.UpdatedAt)
		if err != nil {
			return nil, err
		}
		mentors = append(mentors, mentor)
	}

	return mentors, nil
}

// GetMentorByID finds mentors by ID
func (s *MentorStore) GetMentorByID(ctx context.Context, id int64) (*Mentor, error) {
	query := `SELECT * FROM mentors WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var mentor Mentor
	err := s.db.QueryRowContext(ctx, query, id).Scan(&mentor.ID, &mentor.Userid, &mentor.Name, &mentor.Country, pq.Array(&mentor.Language),
		&mentor.CreatedAt, &mentor.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &mentor, nil
}

// GetMentorsByExpertise finds mentors by expertise
// func (s *MentorStore) GetMentorsByExpertise(ctx context.Context, expertise string) ([]*Mentor, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	rows, err := s.db.QueryContext(ctx, `
// 		SELECT id, userid, name, country, language,created_at, updated_at
// 		FROM mentors
// 		WHERE expertise ILIKE $1
// 		ORDER BY name
// 	`, "%"+expertise+"%")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	mentors := []*Mentor{}
// 	for rows.Next() {
// 		mentor := &Mentor{}
// 		err := rows.Scan(
// 			&mentor.ID, &mentor.Userid, &mentor.Name, &mentor.Country, pq.Array(&mentor.Language),
// 			&mentor.CreatedAt, &mentor.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		mentors = append(mentors, mentor)
// 	}

// 	return mentors, nil
// }

// // GetMentorsByDiscipline finds mentors by discipline
// func (s *MentorStore) GetMentorsByDiscipline(ctx context.Context, discipline string) ([]*Mentor, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	rows, err := s.db.QueryContext(ctx, `
// 		SELECT m.id, m.userid, m.name, m.country, m.language, m.created_at, m.updated_at
// 		FROM mentors m
// 		JOIN Discipline md ON m.id = md.mentor_id
// 		WHERE md.discipline ILIKE $1
// 		ORDER BY m.name
// 	`, "%"+discipline+"%")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	mentors := []*Mentor{}
// 	for rows.Next() {
// 		mentor := &Mentor{}
// 		err := rows.Scan(
// 			&mentor.ID, &mentor.Userid, &mentor.Name, &mentor.Country, pq.Array(&mentor.Language),
// 			&mentor.CreatedAt, &mentor.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		mentors = append(mentors, mentor)
// 	}

// 	return mentors, nil
// }

// UpdateMentor updates an existing mentor
func (s *MentorStore) UpdateMentor(ctx context.Context, mentor *Mentor) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update basic details
	_, err = tx.ExecContext(ctx, `
		UPDATE mentors 
		SET name = $1, country = $2, language = $3, updated_at = NOW()
		WHERE id = $4
	`, mentor.Name, mentor.Country, pq.Array(mentor.Language), mentor.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteMentor removes a mentor
func (s *MentorStore) DeleteMentor(ctx context.Context, mentorID int64) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, `
		DELETE FROM mentors WHERE id = $1
	`, mentorID)
	return err
}
