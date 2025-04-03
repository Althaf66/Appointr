package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	"github.com/go-chi/chi/v5"
)

type RegisterEducationPayload struct {
	Year_from string `json:"year_from" validate:"required"`
	Year_to   string `json:"year_to"`
	Degree    string `json:"degree" validate:"required"`
	Field     string `json:"field"`
	Institute string `json:"institute"`
}

type UpdateEducationPayload struct {
	Year_from *string `json:"year_from"`
	Year_to   *string `json:"year_to"`
	Degree    *string `json:"degree"`
	Field     *string `json:"field"`
	Institute *string `json:"institute"`
}

// createEducationHandler godoc
//
//	@Summary		Create a new education
//	@Description	Create a new education
//	@Tags			education
//	@Accept			json
//	@Produce		json
//	@Param			education	body		RegisterEducationPayload	true	"Education"
//	@Success		200			{object}	error
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/education/create [post]
func (app *application) createEducationHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)

	var payload RegisterEducationPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	education := &store.Education{
		Userid:    user.ID,
		Year_from: payload.Year_from,
		Year_to:   payload.Year_to,
		Degree:    payload.Degree,
		Field:     payload.Field,
		Institute: payload.Institute,
	}

	err = app.store.Education.CreateEducation(r.Context(), education)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	err = JsonResponse(w, http.StatusCreated, education)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getEducationByUserIDHandler godoc
//
//	@Summary		Get all education by user ID
//	@Description	Get all education by user ID
//	@Tags			education
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{array}		store.Education
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/education/u/{id} [get]
func (app *application) getEducationByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	education, err := app.store.Education.GetEducationByUserId(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, education)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getEducationByIDHandler godoc
//
//	@Summary		Get education by ID
//	@Description	Get education by ID
//	@Tags			education
//	@Accept			json
//	@Produce		json
//	@Param			educationid	path		int	true	"Education ID"
//	@Success		200			{object}	store.Education
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/education/{educationid} [get]
func (app *application) getEducationByIDHandler(w http.ResponseWriter, r *http.Request) {
	educationID, err := strconv.ParseInt(chi.URLParam(r, "educationID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	education, err := app.store.Education.GetEducationById(r.Context(), educationID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, education)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// updateEducationHandler godoc
//
// @Summary		Update education
// @Description	Update education
// @Tags			education
// @Accept			json
// @Produce		json
// @Param			educationID	path		int						true	"Education ID"
// @Param			education	body		UpdateEducationPayload	true	"Education"
// @Success		200			{object}	store.Education
// @Failure		404			{object}	error
// @Failure		500			{object}	error
// @Security		ApiKeyAuth
// @Router			/education/{educationID} [patch]
func (app *application) updateEducationHandler(w http.ResponseWriter, r *http.Request) {
	educationID, err := strconv.ParseInt(chi.URLParam(r, "educationID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	education, err := app.store.Education.GetEducationById(r.Context(), educationID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	var payload UpdateEducationPayload
	err = ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Since all fields are optional, we only validate if at least one field is provided
	hasUpdates := payload.Year_from != nil || payload.Year_to != nil ||
		payload.Degree != nil || payload.Field != nil || payload.Institute != nil
	if !hasUpdates {
		app.badRequestResponse(w, r, errors.New("no fields provided for update"))
		return
	}

	// Update only provided fields
	if payload.Year_from != nil {
		education.Year_from = *payload.Year_from
	}
	if payload.Year_to != nil {
		education.Year_to = *payload.Year_to
	}
	if payload.Degree != nil {
		education.Degree = *payload.Degree
	}
	if payload.Field != nil {
		education.Field = *payload.Field
	}
	if payload.Institute != nil {
		education.Institute = *payload.Institute
	}

	err = app.store.Education.UpdateEducation(r.Context(), education)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, education)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// deleteGigHandler godoc
//
//	@Summary		Delete education
//	@Description	Delete education by id
//	@Tags			education
//	@Accept			json
//	@Produce		json
//	@Param			educationID	path		int	true	"Education ID"
//	@Success		200			{object}	nil
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/education/{educationID} [delete]
func (app *application) deleteEducationHandler(w http.ResponseWriter, r *http.Request) {
	educationID, err := strconv.ParseInt(chi.URLParam(r, "educationID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.Education.DeleteEducation(r.Context(), educationID)
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
