package store

import (
	"context"
	"database/sql"
)

type Meetings struct {
	ID          int64   `json:"id"`
	Userid      int64   `json:"userid"`
	Mentorid    int64   `json:"mentorid"`
	Day         string  `json:"day"`
	Date        string  `json:"date"`
	StartTime   string  `json:"start_time"`
	StartPeriod string  `json:"start_period"`
	Isconfirm   bool    `json:"isconfirm"`
	Ispaid      bool    `json:"ispaid"`
	Iscompleted bool    `json:"iscompleted"`
	Amount      float64 `json:"amount"`
	Link        string  `json:"link"`
}

type MeetingsStore struct {
	db *sql.DB
}

func (s *MeetingsStore) CreateMeeting(ctx context.Context, meeting *Meetings) error {
	query := `
		INSERT INTO meetings (userid, mentorid, day, date, start_time, start_period, isconfirm, ispaid, iscompleted, amount, link)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, query,
		meeting.Userid, meeting.Mentorid, meeting.Day, meeting.Date, meeting.StartTime,
		meeting.StartPeriod, meeting.Isconfirm, meeting.Ispaid, meeting.Iscompleted,
		meeting.Amount, meeting.Link).Scan(&meeting.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *MeetingsStore) GetAllMeetings(ctx context.Context, limit, offset int) ([]*Meetings, error) {
	query := `
		SELECT id, userid, mentorid, day, date, start_time, start_period, isconfirm, ispaid, iscompleted, amount, link
		FROM meetings
		ORDER BY id
		LIMIT $1 OFFSET $2`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meetings := []*Meetings{}
	for rows.Next() {
		meeting := &Meetings{}
		err := rows.Scan(
			&meeting.ID, &meeting.Userid, &meeting.Mentorid, &meeting.Day, &meeting.Date,
			&meeting.StartTime, &meeting.StartPeriod, &meeting.Isconfirm, &meeting.Ispaid,
			&meeting.Iscompleted, &meeting.Amount, &meeting.Link,
		)
		if err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

func (s *MeetingsStore) GetMeetingByID(ctx context.Context, id int64) (*Meetings, error) {
	query := `
		SELECT id, userid, mentorid, day, date, start_time, start_period, isconfirm, ispaid, iscompleted, amount, link
		FROM meetings
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	meeting := &Meetings{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&meeting.ID, &meeting.Userid, &meeting.Mentorid, &meeting.Day, &meeting.Date,
		&meeting.StartTime, &meeting.StartPeriod, &meeting.Isconfirm, &meeting.Ispaid,
		&meeting.Iscompleted, &meeting.Amount, &meeting.Link,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return meeting, nil
}

func (s *MeetingsStore) GetMeetingByUserID(ctx context.Context, userid int64) ([]*Meetings, error) {
	query := `
		SELECT id, userid, mentorid, day, date, start_time, start_period, isconfirm, ispaid, iscompleted, amount, link
		FROM meetings
		WHERE userid = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meetings := []*Meetings{}
	for rows.Next() {
		meeting := &Meetings{}
		err := rows.Scan(
			&meeting.ID, &meeting.Userid, &meeting.Mentorid, &meeting.Day, &meeting.Date,
			&meeting.StartTime, &meeting.StartPeriod, &meeting.Isconfirm, &meeting.Ispaid,
			&meeting.Iscompleted, &meeting.Amount, &meeting.Link,
		)
		if err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

func (s *MeetingsStore) GetMeetingMentorNotConfirm(ctx context.Context, mentorID int64) ([]*Meetings, error) {
	query := `
		SELECT id, userid, mentorid, day, date, start_time, start_period, isconfirm, ispaid, iscompleted, amount, link
		FROM meetings
		WHERE mentorid = $1 AND isconfirm = false`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, mentorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meetings := []*Meetings{}
	for rows.Next() {
		meeting := &Meetings{}
		err := rows.Scan(
			&meeting.ID, &meeting.Userid, &meeting.Mentorid, &meeting.Day, &meeting.Date,
			&meeting.StartTime, &meeting.StartPeriod, &meeting.Isconfirm, &meeting.Ispaid,
			&meeting.Iscompleted, &meeting.Amount, &meeting.Link,
		)
		if err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

func (s *MeetingsStore) GetMeetingUserNotPaid(ctx context.Context, userID int64) ([]*Meetings, error) {
	query := `
		SELECT id, userid, mentorid, day, date, start_time, start_period, isconfirm, ispaid, iscompleted, amount, link
		FROM meetings
		WHERE userid = $1 AND isconfirm = true AND ispaid = false`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meetings := []*Meetings{}
	for rows.Next() {
		meeting := &Meetings{}
		err := rows.Scan(
			&meeting.ID, &meeting.Userid, &meeting.Mentorid, &meeting.Day, &meeting.Date,
			&meeting.StartTime, &meeting.StartPeriod, &meeting.Isconfirm, &meeting.Ispaid,
			&meeting.Iscompleted, &meeting.Amount, &meeting.Link,
		)
		if err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

func (s *MeetingsStore) GetMeetingUserNotCompleted(ctx context.Context, userID int64) ([]*Meetings, error) {
	query := `
		SELECT id, userid, mentorid, day, date, start_time, start_period, isconfirm, ispaid, iscompleted, amount, link
		FROM meetings
		WHERE userid = $1 AND isconfirm = true AND ispaid = true AND iscompleted = false`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meetings := []*Meetings{}
	for rows.Next() {
		meeting := &Meetings{}
		err := rows.Scan(
			&meeting.ID, &meeting.Userid, &meeting.Mentorid, &meeting.Day, &meeting.Date,
			&meeting.StartTime, &meeting.StartPeriod, &meeting.Isconfirm, &meeting.Ispaid,
			&meeting.Iscompleted, &meeting.Amount, &meeting.Link,
		)
		if err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

func (s *MeetingsStore) GetMeetingMentorNotCompleted(ctx context.Context, mentorID int64) ([]*Meetings, error) {
	query := `
		SELECT id, userid, mentorid, day, date, start_time, start_period, isconfirm, ispaid, iscompleted, amount, link
		FROM meetings
		WHERE mentorid = $1 AND isconfirm = true AND ispaid = true AND iscompleted = false`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, mentorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meetings := []*Meetings{}
	for rows.Next() {
		meeting := &Meetings{}
		err := rows.Scan(
			&meeting.ID, &meeting.Userid, &meeting.Mentorid, &meeting.Day, &meeting.Date,
			&meeting.StartTime, &meeting.StartPeriod, &meeting.Isconfirm, &meeting.Ispaid,
			&meeting.Iscompleted, &meeting.Amount, &meeting.Link,
		)
		if err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}

	return meetings, nil
}

func (s *MeetingsStore) UpdateMeetingConfirm(ctx context.Context, meetingID int64) error {
	query := `
		UPDATE meetings 
		SET isconfirm = true
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, query, meetingID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return tx.Commit()
}

func (s *MeetingsStore) UpdateMeetingPaid(ctx context.Context, meetingID int64) error {
	query := `
		UPDATE meetings 
		SET ispaid = true
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, query, meetingID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return tx.Commit()
}

func (s *MeetingsStore) UpdateMeetingCompleted(ctx context.Context, meetingID int64) error {
	query := `
		UPDATE meetings 
		SET iscompleted = true
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, query, meetingID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return tx.Commit()
}

func (s *MeetingsStore) DeleteMeeting(ctx context.Context, meetingID int64) error {
	query := `
		DELETE FROM meetings 
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, query, meetingID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return tx.Commit()
}
