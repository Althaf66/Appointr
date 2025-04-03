package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type BookingSlot struct {
    ID          int64      `json:"id"`
    UserID      int64     `json:"userid"`
    Days        []string `json:"days"`
    StartTime   string   `json:"start_time"`
    StartPeriod string   `json:"start_period"`
    EndTime     string   `json:"end_time"`
    EndPeriod   string   `json:"end_period"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type BookingStore struct {
	db *sql.DB
}

func (s *BookingStore) CreateBookingSlot(ctx context.Context, slot *BookingSlot) error {
	query := `INSERT INTO bookingslots (userid, days, start_time, start_period, end_time, end_period)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.db.QueryRowContext(ctx, query, slot.UserID, pq.Array(slot.Days), slot.StartTime,
		slot.StartPeriod, slot.EndTime, slot.EndPeriod).Scan(&slot.ID,&slot.CreatedAt,&slot.UpdatedAt)
	if err != nil {
		return err
	}
	// createdSlots := make([]BookingSlot, 0, len(slots))
    // for _, slot := range slots {
    //     var id int
    //     var createdAt, updatedAt time.Time
    //     err = s.db.QueryRow(
    //         userID,
    //         slot.Days,
    //         slot.StartTime,
    //         slot.StartPeriod,
    //         slot.EndTime,
    //         slot.EndPeriod,
    //     ).Scan(&id, &createdAt, &updatedAt)
		// createdSlots = append(createdSlots, BookingSlot{
        //     ID:          id,
        //     UserID:      userID,
        //     Days:        slot.Days,
        //     StartTime:   slot.StartTime,
        //     StartPeriod: slot.StartPeriod,
        //     EndTime:     slot.EndTime,
        //     EndPeriod:   slot.EndPeriod,
        //     CreatedAt:   createdAt,
        //     UpdatedAt:   updatedAt,
        // })
    // }		
	return tx.Commit()
}