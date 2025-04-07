package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	chi "github.com/go-chi/chi/v5"
)

type RegisterMeetingPayload struct {
	Mentorid    int64   `json:"mentorid"`
	Day         string  `json:"day"`
	Date        string  `json:"date"`
	StartTime   string  `json:"start_time"`
	StartPeriod string  `json:"start_period"`
	Amount      float64 `json:"amount"`
}

// createMeetingHandler godoc
//
//	@Summary		Create a new meeting
//	@Description	Create a new meeting
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Param			meeting	body		RegisterMeetingPayload	true	"Meeting"
//	@Success		201		{object}	store.Meetings
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings/create [post]
func (app *application) createMeetingHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)

	var payload RegisterMeetingPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	meeting := &store.Meetings{
		Userid:      user.ID,
		Mentorid:    payload.Mentorid,
		Day:         payload.Day,
		Date:        payload.Date,
		StartTime:   payload.StartTime,
		StartPeriod: payload.StartPeriod,
		Isconfirm:   false,
		Ispaid:      false,
		Iscompleted: false,
		Amount:      payload.Amount,
	}

	err = app.store.Meetings.CreateMeeting(r.Context(), meeting)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	err = JsonResponse(w, http.StatusCreated, meeting)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getAllMeetingsHandler godoc
//
//	@Summary		Get all meetings
//	@Description	Get all meetings
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]store.Meetings
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings [get]
func (app *application) getAllMeetingsHandler(w http.ResponseWriter, r *http.Request) {
	meetings, err := app.store.Meetings.GetAllMeetings(r.Context(), 10, 0)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, meetings)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getMeetingByUserIDHandler godoc
//
//	@Summary		Get meetings by user ID
//	@Description	Get meetings for a specific user
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64	true	"User ID"
//	@Success		200		{object}	[]store.Meetings
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings/u/{userID} [get]
func (app *application) getMeetingByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	meetings, err := app.store.Meetings.GetMeetingByUserID(r.Context(), userID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, meetings)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getMeetingMentorNotConfirmHandler godoc
//
//	@Summary		Get unconfirmed meetings by mentor ID
//	@Description	Get meetings where isconfirm is false for a specific mentor
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Param			mentorID	path		int64	true	"Mentor ID"
//	@Success		200			{object}	[]store.Meetings
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings/mentor-not-confirm/{mentorID} [get]
func (app *application) getMeetingMentorNotConfirmHandler(w http.ResponseWriter, r *http.Request) {
	mentorID, err := strconv.ParseInt(chi.URLParam(r, "mentorID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	meetings, err := app.store.Meetings.GetMeetingMentorNotConfirm(r.Context(), mentorID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, meetings)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getMeetingUserNotPaidHandler godoc
//
//	@Summary		Get unpaid meetings by user ID
//	@Description	Get meetings where ispaid is false for a specific user
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64	true	"User ID"
//	@Success		200		{object}	[]store.Meetings
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings/user-not-paid/{userID} [get]
func (app *application) getMeetingUserNotPaidHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	meetings, err := app.store.Meetings.GetMeetingUserNotPaid(r.Context(), userID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, meetings)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getMeetingUserNotCompletedHandler godoc
//
//	@Summary		Get uncompleted meetings by user ID
//	@Description	Get meetings where iscompleted is false but confirmed and paid for a specific user
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64	true	"User ID"
//	@Success		200		{object}	[]store.Meetings
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings/user-not-completed/{userID} [get]
func (app *application) getMeetingUserNotCompletedHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	meetings, err := app.store.Meetings.GetMeetingUserNotCompleted(r.Context(), userID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, meetings)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getMeetingMentorNotCompletedHandler godoc
//
//	@Summary		Get uncompleted meetings by mentor ID
//	@Description	Get meetings where iscompleted is false but confirmed and paid for a specific mentor
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Param			mentorID	path		int64	true	"Mentor ID"
//	@Success		200			{object}	[]store.Meetings
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings/mentor-not-completed/{mentorID} [get]
func (app *application) getMeetingMentorNotCompletedHandler(w http.ResponseWriter, r *http.Request) {
	mentorID, err := strconv.ParseInt(chi.URLParam(r, "mentorID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	meetings, err := app.store.Meetings.GetMeetingMentorNotCompleted(r.Context(), mentorID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, meetings)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// updateMeetingConfirmHandler godoc
//
//	@Summary		Confirm a meeting
//	@Description	Update meeting confirmation status to true
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Param			meetingID	path		int64	true	"Meeting ID"
//	@Success		200			{object}	store.Meetings
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings/confirm/{meetingID} [put]
func (app *application) updateMeetingConfirmHandler(w http.ResponseWriter, r *http.Request) {
	meetingID, err := strconv.ParseInt(chi.URLParam(r, "meetingID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	meeting, err := app.store.Meetings.GetMeetingByID(r.Context(), meetingID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.store.Meetings.UpdateMeetingConfirm(r.Context(), meetingID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	meeting.Isconfirm = true
	err = JsonResponse(w, http.StatusOK, meeting)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// updateMeetingPaidHandler godoc
//
//	@Summary		Mark meeting as paid
//	@Description	Update meeting paid status to true
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Param			meetingID	path		int64	true	"Meeting ID"
//	@Success		200			{object}	store.Meetings
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings/paid/{meetingID} [put]
func (app *application) updateMeetingPaidHandler(w http.ResponseWriter, r *http.Request) {
	meetingID, err := strconv.ParseInt(chi.URLParam(r, "meetingID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	meeting, err := app.store.Meetings.GetMeetingByID(r.Context(), meetingID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.store.Meetings.UpdateMeetingPaid(r.Context(), meetingID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	meeting.Ispaid = true
	err = JsonResponse(w, http.StatusOK, meeting)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// updateMeetingCompletedHandler godoc
//
//	@Summary		Mark meeting as completed
//	@Description	Update meeting completed status to true
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Param			meetingID	path		int64	true	"Meeting ID"
//	@Success		200			{object}	store.Meetings
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings/completed/{meetingID} [put]
func (app *application) updateMeetingCompletedHandler(w http.ResponseWriter, r *http.Request) {
	meetingID, err := strconv.ParseInt(chi.URLParam(r, "meetingID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	meeting, err := app.store.Meetings.GetMeetingByID(r.Context(), meetingID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.store.Meetings.UpdateMeetingCompleted(r.Context(), meetingID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	meeting.Iscompleted = true
	err = JsonResponse(w, http.StatusOK, meeting)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// deleteMeetingHandler godoc
//
//	@Summary		Delete meeting
//	@Description	Delete a specific meeting
//	@Tags			meetings
//	@Accept			json
//	@Produce		json
//	@Param			meetingID	path		int64	true	"Meeting ID"
//	@Success		200			{object}	nil
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/meetings/{meetingID} [delete]
func (app *application) deleteMeetingHandler(w http.ResponseWriter, r *http.Request) {
	meetingID, err := strconv.ParseInt(chi.URLParam(r, "meetingID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.Meetings.DeleteMeeting(r.Context(), meetingID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
