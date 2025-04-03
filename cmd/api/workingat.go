package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	"github.com/go-chi/chi/v5"
)

type RegisterWorkingAtPayload struct {
	Title     string `json:"title" validate:"required"`
	Company   string `json:"company"`
	TotalYear int64  `json:"totalyear"`
	Month     int64  `json:"month"`
	Linkedin  string `json:"linkedin"`
	Github    string `json:"github"`
	Instagram string `json:"instagram"`
}

type UpdateWorkingAtPayload struct {
	Title     *string `json:"title"`
	Company   *string `json:"company"`
	TotalYear *int64  `json:"totalyear"`
	Month     *int64  `json:"month"`
	Linkedin  *string `json:"linkedin"`
	Github    *string `json:"github"`
	Instagram *string `json:"instagram"`
}

// createWorkingAtHandler godoc
//
//	@Summary		Create a new workingat
//	@Description	Create a new workingat
//	@Tags			workingat
//	@Accept			json
//	@Produce		json
//	@Param			workingat	body		RegisterWorkingAtPayload	true	"WorkingAt"
//	@Success		200			{object}	error
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/workingat/create [post]
func (app *application) createWorkingAtHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)

	var payload RegisterWorkingAtPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workingat := &store.WorkingAt{
		Userid:    user.ID,
		Title:     payload.Title,
		Company:   payload.Company,
		TotalYear: payload.TotalYear,
		Month:     payload.Month,
		Linkedin:  payload.Linkedin,
		Github:    payload.Github,
		Instagram: payload.Instagram,
	}

	err = app.store.WorkingAt.CreateWorkingAt(r.Context(), workingat)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	err = JsonResponse(w, http.StatusCreated, workingat)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getWorkingAtByUserIDHandler godoc
//
//	@Summary		Get all workingat by user ID
//	@Description	Get all workingat by user ID
//	@Tags			workingat
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{array}		store.WorkingAt
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/workingat/u/{id} [get]
func (app *application) getWorkingAtByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workingat, err := app.store.WorkingAt.GetWorkingAtByUserId(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, workingat)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getWorkingAtByIDHandler godoc
//
//	@Summary		Get workingat by ID
//	@Description	Get workingat by ID
//	@Tags			workingat
//	@Accept			json
//	@Produce		json
//	@Param			workingatID	path		int	true	"WorkingAt ID"
//	@Success		200			{object}	store.WorkingAt
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/workingat/{workingatID} [get]
func (app *application) getWorkingAtByIDHandler(w http.ResponseWriter, r *http.Request) {
	workingatID, err := strconv.ParseInt(chi.URLParam(r, "workingatID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workingat, err := app.store.WorkingAt.GetWorkingAtById(r.Context(), workingatID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, workingat)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// updateWorkingAtHandler godoc
//
//	@Summary		Update workingat
//	@Description	Update workingat
//	@Tags			workingat
//	@Accept			json
//	@Produce		json
//	@Param			workingatID	path		int						true	"WorkingAt ID"
//	@Param			workingat	body		UpdateWorkingAtPayload	true	"WorkingAt"
//	@Success		200			{object}	store.WorkingAt
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/workingat/{workingatID} [patch]
func (app *application) updateWorkingAtHandler(w http.ResponseWriter, r *http.Request) {
	workingatID, err := strconv.ParseInt(chi.URLParam(r, "workingatID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workingat, err := app.store.WorkingAt.GetWorkingAtById(r.Context(), workingatID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	var payload UpdateWorkingAtPayload
	err = ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Since all companys are optional, we only validate if at least one company is provided
	hasUpdates := payload.TotalYear != nil || payload.Month != nil ||
		payload.Title != nil || payload.Company != nil
	if !hasUpdates {
		app.badRequestResponse(w, r, errors.New("no companys provided for update"))
		return
	}

	// Update only provided companys
	if payload.Title != nil {
		workingat.Title = *payload.Title
	}
	if payload.Company != nil {
		workingat.Company = *payload.Company
	}
	if payload.TotalYear != nil {
		workingat.TotalYear = *payload.TotalYear
	}
	if payload.Month != nil {
		workingat.Month = *payload.Month
	}
	if payload.Linkedin != nil {
		workingat.Linkedin = *payload.Linkedin
	}
	if payload.Github != nil {
		workingat.Github = *payload.Github
	}
	if payload.Instagram != nil {
		workingat.Instagram = *payload.Instagram
	}

	err = app.store.WorkingAt.UpdateWorkingAt(r.Context(), workingat)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, workingat)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// deleteGigHandler godoc
//
//	@Summary		Delete workingat
//	@Description	Delete workingat by id
//	@Tags			workingat
//	@Accept			json
//	@Produce		json
//	@Param			workingatID	path		int	true	"WorkingAt ID"
//	@Success		200			{object}	nil
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/workingat/{workingatID} [delete]
func (app *application) deleteWorkingAtHandler(w http.ResponseWriter, r *http.Request) {
	workingatID, err := strconv.ParseInt(chi.URLParam(r, "workingatID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.WorkingAt.DeleteWorkingAt(r.Context(), workingatID)
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
