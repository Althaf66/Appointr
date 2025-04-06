package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var ErrNotFound = errors.New("resource not found")

type Mentor struct {
	ID          int64         `json:"id"`
	Userid      int64         `json:"userid"`
	Name        string        `json:"name"`
	Country     string        `json:"country"`
	Language    []string      `json:"language"`
	Gigs        []Gig         `json:"gigs"`
	Education   []Education   `json:"education"`
	Experience  []Experience  `json:"experience"`
	WorkingAt   *WorkingAt    `json:"workingat"`
	BookingSlot []BookingSlot `json:"bookingslots"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
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

// GetMentorByUserID finds mentors by userid
func (s *MentorStore) GetMentorByUserID(ctx context.Context, userid int64) (*Mentor, error) {
	mentor := &Mentor{}
	err := s.db.QueryRow(`
        SELECT id, userid, name, country, language, created_at, updated_at
        FROM mentors
        WHERE userid = $1`, userid).Scan(
		&mentor.ID,
		&mentor.Userid,
		&mentor.Name,
		&mentor.Country,
		pq.Array(&mentor.Language),
		&mentor.CreatedAt,
		&mentor.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mentor with userid %d not found", userid)
		}
		return nil, err
	}

	// Fetch gigs
	gigsRows, err := s.db.Query(`
        SELECT id, title, amount,description, expertise, discipline, created_at, updated_at
        FROM gigs
        WHERE userid = $1`, userid)
	if err != nil {
		return nil, err
	}
	defer gigsRows.Close()

	for gigsRows.Next() {
		gig := Gig{}
		err = gigsRows.Scan(
			&gig.ID,
			&gig.Title,
			&gig.Amount,
			&gig.Description,
			&gig.Expertise,
			pq.Array(&gig.Discipline),
			&gig.CreatedAt,
			&gig.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		mentor.Gigs = append(mentor.Gigs, gig)
	}

	// Fetch education
	eduRows, err := s.db.Query(`
        SELECT id, year_from, year_to, degree, field, institute
        FROM education
        WHERE userid = $1`, userid)
	if err != nil {
		return nil, err
	}
	defer eduRows.Close()

	for eduRows.Next() {
		edu := Education{}
		err = eduRows.Scan(&edu.ID, &edu.Year_from, &edu.Year_to, &edu.Degree, &edu.Field, &edu.Institute)
		if err != nil {
			return nil, err
		}
		mentor.Education = append(mentor.Education, edu)
	}

	// Fetch experience
	expRows, err := s.db.Query(`
        SELECT id, year_from, year_to, title, company, description
        FROM experience
        WHERE userid = $1`, userid)
	if err != nil {
		return nil, err
	}
	defer expRows.Close()

	for expRows.Next() {
		exp := Experience{}
		err = expRows.Scan(&exp.ID, &exp.Year_from, &exp.Year_to, &exp.Title, &exp.Company, &exp.Description)
		if err != nil {
			return nil, err
		}
		mentor.Experience = append(mentor.Experience, exp)
	}

	// Fetch working at (single record)
	workingAt := WorkingAt{}
	err = s.db.QueryRow(`
        SELECT id, title, company, totalyear, month
        FROM workingat
        WHERE userid = $1`, userid).Scan(
		&workingAt.ID,
		&workingAt.Title,
		&workingAt.Company,
		&workingAt.TotalYear,
		&workingAt.Month,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == nil {
		mentor.WorkingAt = &workingAt
	}

	// Fetch BookingSlot
	bookingRows, err := s.db.Query(`
        SELECT id, days, start_time, start_period, end_time, end_period
        FROM bookingslots
        WHERE userid = $1`, userid)
	if err != nil {
		return nil, err
	}
	defer bookingRows.Close()

	for bookingRows.Next() {
		slot := BookingSlot{}
		err = bookingRows.Scan(&slot.ID, pq.Array(&slot.Days), &slot.StartTime, &slot.StartPeriod, &slot.EndTime, &slot.EndPeriod)
		if err != nil {
			return nil, err
		}
		mentor.BookingSlot = append(mentor.BookingSlot, slot)
	}

	return mentor, nil

}

func (s *MentorStore) GetMentorsByExpertise(ctx context.Context, expertise string) ([]*Mentor, error) {
	// First get all userIDs with matching expertise
	rows, err := s.db.Query(`
        SELECT DISTINCT m.id, m.userid, m.name, m.country, m.language, m.created_at, m.updated_at
        FROM mentors m
        JOIN gigs g ON m.userid = g.userid
        WHERE g.expertise = $1`, expertise)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mentors []*Mentor
	mentorMap := make(map[int64]*Mentor)

	// Collect basic mentor info
	for rows.Next() {
		mentor := &Mentor{}
		err = rows.Scan(
			&mentor.ID,
			&mentor.Userid,
			&mentor.Name,
			&mentor.Country,
			pq.Array(&mentor.Language),
			&mentor.CreatedAt,
			&mentor.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		mentors = append(mentors, mentor)
		mentorMap[mentor.Userid] = mentor
	}

	if len(mentors) == 0 {
		return nil, fmt.Errorf("no mentors found with expertise: %s", expertise)
	}

	// Get all userIDs for batch queries
	var userIDs []int64
	for _, m := range mentors {
		userIDs = append(userIDs, m.Userid)
	}

	// Fetch gigs
	gigsRows, err := s.db.Query(`
        SELECT id, userid, title, amount,description, expertise, discipline, created_at, updated_at
        FROM gigs
        WHERE userid = ANY($1) AND expertise = $2`, pq.Array(userIDs), expertise)
	if err != nil {
		return nil, err
	}
	defer gigsRows.Close()

	for gigsRows.Next() {
		gig := Gig{}
		var userID int64
		err = gigsRows.Scan(
			&gig.ID,
			&userID,
			&gig.Title,
			&gig.Amount,
			&gig.Description,
			&gig.Expertise,
			pq.Array(&gig.Discipline),
			&gig.CreatedAt,
			&gig.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if mentor, ok := mentorMap[userID]; ok {
			mentor.Gigs = append(mentor.Gigs, gig)
		}
	}

	// Fetch education
	eduRows, err := s.db.Query(`
        SELECT id, userid, year_from, year_to, degree, field, institute
        FROM education
        WHERE userid = ANY($1)`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer eduRows.Close()

	for eduRows.Next() {
		edu := Education{}
		var userID int64
		err = eduRows.Scan(&edu.ID, &userID, &edu.Year_from, &edu.Year_to, &edu.Degree, &edu.Field, &edu.Institute)
		if err != nil {
			return nil, err
		}
		if mentor, ok := mentorMap[userID]; ok {
			mentor.Education = append(mentor.Education, edu)
		}
	}

	// Fetch experience
	expRows, err := s.db.Query(`
        SELECT id, userid, year_from, year_to, title, company, description
        FROM experience
        WHERE userid = ANY($1)`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer expRows.Close()

	for expRows.Next() {
		exp := Experience{}
		var userID int64
		err = expRows.Scan(&exp.ID, &userID, &exp.Year_from, &exp.Year_to, &exp.Title, &exp.Company, &exp.Description)
		if err != nil {
			return nil, err
		}
		if mentor, ok := mentorMap[userID]; ok {
			mentor.Experience = append(mentor.Experience, exp)
		}
	}

	// Fetch working at
	workRows, err := s.db.Query(`
        SELECT id, userid, title, company, totalyear, month
        FROM workingat
        WHERE userid = ANY($1)`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer workRows.Close()

	for workRows.Next() {
		work := WorkingAt{}
		var userID int64
		err = workRows.Scan(&work.ID, &userID, &work.Title, &work.Company, &work.TotalYear, &work.Month)
		if err != nil {
			return nil, err
		}
		if mentor, ok := mentorMap[userID]; ok {
			mentor.WorkingAt = &work
		}
	}

	// Fetch booking slot
	bookingRows, err := s.db.Query(`
        SELECT id,userid, days, start_time, start_period, end_time, end_period
        FROM bookingslots
        WHERE userid = ANY($1)`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer bookingRows.Close()

	for bookingRows.Next() {
		slot := BookingSlot{}
		var userID int64
		err = bookingRows.Scan(&slot.ID, &userID, pq.Array(&slot.Days), &slot.StartTime, &slot.StartPeriod, &slot.EndTime, &slot.EndPeriod)
		if err != nil {
			return nil, err
		}
		if mentor, ok := mentorMap[userID]; ok {
			mentor.BookingSlot = append(mentor.BookingSlot, slot)
		}
	}

	return mentors, nil
}

func (s *MentorStore) GetMentorsByDiscipline(ctx context.Context, discipline string) ([]*Mentor, error) {
	// First get all mentors with matching discipline using array contains operator
	rows, err := s.db.Query(`
        SELECT DISTINCT m.id, m.userid, m.name, m.country, m.language, m.created_at, m.updated_at
        FROM mentors m
        JOIN gigs g ON m.userid = g.userid
        WHERE $1 = ANY(g.discipline)`, discipline)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mentors []*Mentor
	mentorMap := make(map[int64]*Mentor)

	// Collect basic mentor info
	for rows.Next() {
		mentor := &Mentor{}
		err = rows.Scan(
			&mentor.ID,
			&mentor.Userid,
			&mentor.Name,
			&mentor.Country,
			pq.Array(&mentor.Language),
			&mentor.CreatedAt,
			&mentor.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		mentors = append(mentors, mentor)
		mentorMap[mentor.Userid] = mentor
	}

	if len(mentors) == 0 {
		return nil, fmt.Errorf("no mentors found with discipline: %s", discipline)
	}

	// Get all userIDs for batch queries
	var userIDs []int64
	for _, m := range mentors {
		userIDs = append(userIDs, m.Userid)
	}

	// Fetch gigs (only those with matching discipline)
	gigsRows, err := s.db.Query(`
        SELECT id, userid, title, amount,description, expertise, discipline, created_at, updated_at
        FROM gigs
        WHERE userid = ANY($1) AND $2 = ANY(discipline)`, pq.Array(userIDs), discipline)
	if err != nil {
		return nil, err
	}
	defer gigsRows.Close()

	for gigsRows.Next() {
		gig := Gig{}
		var userID int64
		err = gigsRows.Scan(
			&gig.ID,
			&userID,
			&gig.Title,
			&gig.Amount,
			&gig.Description,
			&gig.Expertise,
			pq.Array(&gig.Discipline),
			&gig.CreatedAt,
			&gig.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if mentor, ok := mentorMap[userID]; ok {
			mentor.Gigs = append(mentor.Gigs, gig)
		}
	}

	// Fetch education
	eduRows, err := s.db.Query(`
        SELECT id, userid, year_from, year_to, degree, field, institute
        FROM education
        WHERE userid = ANY($1)`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer eduRows.Close()

	for eduRows.Next() {
		edu := Education{}
		var userID int64
		err = eduRows.Scan(&edu.ID, &userID, &edu.Year_from, &edu.Year_to, &edu.Degree, &edu.Field, &edu.Institute)
		if err != nil {
			return nil, err
		}
		if mentor, ok := mentorMap[userID]; ok {
			mentor.Education = append(mentor.Education, edu)
		}
	}

	// Fetch experience
	expRows, err := s.db.Query(`
        SELECT id, userid, year_from, year_to, title, company, description
        FROM experience
        WHERE userid = ANY($1)`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer expRows.Close()

	for expRows.Next() {
		exp := Experience{}
		var userID int64
		err = expRows.Scan(&exp.ID, &userID, &exp.Year_from, &exp.Year_to, &exp.Title, &exp.Company, &exp.Description)
		if err != nil {
			return nil, err
		}
		if mentor, ok := mentorMap[userID]; ok {
			mentor.Experience = append(mentor.Experience, exp)
		}
	}

	// Fetch working at
	workRows, err := s.db.Query(`
        SELECT id, userid, title, company, totalyear, month
        FROM workingat
        WHERE userid = ANY($1)`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer workRows.Close()

	for workRows.Next() {
		work := WorkingAt{}
		var userID int64
		err = workRows.Scan(&work.ID, &userID, &work.Title, &work.Company, &work.TotalYear, &work.Month)
		if err != nil {
			return nil, err
		}
		if mentor, ok := mentorMap[userID]; ok {
			mentor.WorkingAt = &work
		}
	}

	// Fetch bookingslot
	bookingRows, err := s.db.Query(`
        SELECT id, userid, days, start_time, start_period, end_time, end_period
        FROM bookingslots
        WHERE userid = ANY($1)`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer bookingRows.Close()

	for bookingRows.Next() {
		slot := BookingSlot{}
		var userID int64
		err = bookingRows.Scan(&slot.ID, &userID, pq.Array(&slot.Days), &slot.StartTime, &slot.StartPeriod, &slot.EndTime, &slot.EndPeriod)
		if err != nil {
			return nil, err
		}
		if mentor, ok := mentorMap[userID]; ok {
			mentor.BookingSlot = append(mentor.BookingSlot, slot)
		}
	}

	return mentors, nil
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
