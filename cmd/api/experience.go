package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	"github.com/go-chi/chi/v5"
)

type RegisterExperiencePayload struct {
	Year_from   string `json:"year_from" validate:"required"`
	Year_to     string `json:"year_to"`
	Title       string `json:"title" validate:"required"`
	Company     string `json:"company"`
	Description string `json:"description"`
}

type UpdateExperiencePayload struct {
	Year_from   *string `json:"year_from"`
	Year_to     *string `json:"year_to"`
	Title       *string `json:"title"`
	Company     *string `json:"company"`
	Description *string `json:"description"`
}

// createExperienceHandler godoc
//
//	@Summary		Create a new experience
//	@Description	Create a new experience
//	@Tags			experience
//	@Accept			json
//	@Produce		json
//	@Param			experience	body		RegisterExperiencePayload	true	"Experience"
//	@Success		200			{object}	error
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/experience/create [post]
func (app *application) createExperienceHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)

	var payload RegisterExperiencePayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	experience := &store.Experience{
		Userid:      user.ID,
		Year_from:   payload.Year_from,
		Year_to:     payload.Year_to,
		Title:       payload.Title,
		Company:     payload.Company,
		Description: payload.Description,
	}

	err = app.store.Experience.CreateExperience(r.Context(), experience)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	err = JsonResponse(w, http.StatusCreated, experience)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getExperienceByUserIDHandler godoc
//
//	@Summary		Get all experience by user ID
//	@Description	Get all experience by user ID
//	@Tags			experience
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{array}		store.Experience
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/experience/u/{id} [get]
func (app *application) getExperienceByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	experience, err := app.store.Experience.GetExperienceByUserId(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, experience)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getExperienceByIDHandler godoc
//
//	@Summary		Get experience by ID
//	@Description	Get experience by ID
//	@Tags			experience
//	@Accept			json
//	@Produce		json
//	@Param			experienceid	path		int	true	"Experience ID"
//	@Success		200				{object}	store.Experience
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/experience/{experienceid} [get]
func (app *application) getExperienceByIDHandler(w http.ResponseWriter, r *http.Request) {
	experienceID, err := strconv.ParseInt(chi.URLParam(r, "experienceID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	experience, err := app.store.Experience.GetExperienceById(r.Context(), experienceID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, experience)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// updateExperienceHandler godoc
//
// @Summary		Update experience
// @Description	Update experience
// @Tags			experience
// @Accept			json
// @Produce		json
// @Param			experienceID	path		int						true	"Experience ID"
// @Param			experience		body		UpdateExperiencePayload	true	"Experience"
// @Success		200				{object}	store.Experience
// @Failure		404				{object}	error
// @Failure		500				{object}	error
// @Security		ApiKeyAuth
// @Router			/experience/{experienceID} [patch]
func (app *application) updateExperienceHandler(w http.ResponseWriter, r *http.Request) {
	experienceID, err := strconv.ParseInt(chi.URLParam(r, "experienceID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	experience, err := app.store.Experience.GetExperienceById(r.Context(), experienceID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	var payload UpdateExperiencePayload
	err = ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Since all companys are optional, we only validate if at least one company is provided
	hasUpdates := payload.Year_from != nil || payload.Year_to != nil ||
		payload.Title != nil || payload.Company != nil || payload.Description != nil
	if !hasUpdates {
		app.badRequestResponse(w, r, errors.New("no companys provided for update"))
		return
	}

	// Update only provided companys
	if payload.Year_from != nil {
		experience.Year_from = *payload.Year_from
	}
	if payload.Year_to != nil {
		experience.Year_to = *payload.Year_to
	}
	if payload.Title != nil {
		experience.Title = *payload.Title
	}
	if payload.Company != nil {
		experience.Company = *payload.Company
	}
	if payload.Description != nil {
		experience.Description = *payload.Description
	}

	err = app.store.Experience.UpdateExperience(r.Context(), experience)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, experience)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// deleteGigHandler godoc
//
//	@Summary		Delete experience
//	@Description	Delete experience by id
//	@Tags			experience
//	@Accept			json
//	@Produce		json
//	@Param			experienceID	path		int	true	"Experience ID"
//	@Success		200				{object}	nil
//	@Failure		400				{object}	error
//	@Failure		401				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/experience/{experienceID} [delete]
func (app *application) deleteExperienceHandler(w http.ResponseWriter, r *http.Request) {
	experienceID, err := strconv.ParseInt(chi.URLParam(r, "experienceID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.Experience.DeleteExperience(r.Context(), experienceID)
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
