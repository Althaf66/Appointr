package main

import (
	"net/http"

	"github.com/Althaf66/Appointr/internal/store"
)

type RegisterBookingSlotPayload struct {
	Days        []string `json:"days"`
	StartTime   string   `json:"start_time"`
	StartPeriod string   `json:"start_period"`
	EndTime     string   `json:"end_time"`
	EndPeriod   string   `json:"end_period"`
}

// createBookingSlotHandler godoc
//
//	@Summary		Create a new booking slot
//	@Description	Create a new booking slot
//	@Tags			booking
//	@Accept			json
//	@Produce		json
//	@Param			booking	slot		body	RegisterBookingSlotPayload	true	"Booking Slot"
//	@Success		200		{object}	error
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/bookingslots/create [post]
func (app *application) createBookingSlotHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)

	var payload RegisterBookingSlotPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	bookingslot := &store.BookingSlot{
		UserID:      user.ID,
		Days:        payload.Days,
		StartTime:   payload.StartTime,
		StartPeriod: payload.StartPeriod,
		EndTime:     payload.EndTime,
		EndPeriod:   payload.EndPeriod,
	}

	err = app.store.BookingSlot.CreateBookingSlot(r.Context(), bookingslot)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	err = JsonResponse(w, http.StatusCreated, bookingslot)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
