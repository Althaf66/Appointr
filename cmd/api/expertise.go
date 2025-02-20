package main

import (
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	chi "github.com/go-chi/chi/v5"
)

type RegisterExpertisePayload struct {
	Name     string `json:"name" `
	Icon_svg string `json:"icon_svg"`
}

type UpdateExpertisePayload struct {
	Name     *string `json:"name" `
	Icon_svg *string `json:"icon_svg"`
}

// createExpertiseHandler godoc
//
//	@Summary		Creates a new expertise
//	@Description	creates a new expertise field
//	@Tags			expertise
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterExpertisePayload	true	"expertise"
//	@Success		201		{object}	int64						"Expertise registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/expertise/create [post]
func (app *application) createExpertiseHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterExpertisePayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	expertise := &store.Expertise{
		Name:     payload.Name,
		Icon_svg: payload.Icon_svg,
	}

	err = app.store.Expertise.Create(r.Context(), expertise)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	expertiseID := expertise.ID
	err = JsonResponse(w, http.StatusCreated, expertiseID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getExpertiseHandlerByID godoc
//
//	@Summary		Fetches expertise
//	@Description	Fetches expertise by ID
//	@Tags			expertise
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Expertise ID"
//	@Success		200	{object}	store.Expertise
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/expertise/{id} [get]
func (app *application) getExpertiseHandlerByID(w http.ResponseWriter, r *http.Request) {
	expertiseID, err := strconv.ParseInt(chi.URLParam(r, "expertiseID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	expertise, err := app.store.Expertise.GetByID(r.Context(), expertiseID)
	if err != nil {
		switch err {
		case store.ErrExpertiseNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := JsonResponse(w, http.StatusOK, expertise); err != nil {
		app.internalServerError(w, r, err)
	}
}

// deleteExpertiseHandler godoc
//
//	@Summary		Deletes a expertise field
//	@Description	Delete a expertise field by ID
//	@Tags			expertise
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Discipline ID"
//	@Success		204	{object}	string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/expertise/{id} [delete]
func (app *application) deleteExpertiseHandler(w http.ResponseWriter, r *http.Request) {
	expertiseID, err := strconv.ParseInt(chi.URLParam(r, "expertiseID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.Expertise.Delete(r.Context(), expertiseID)
	if err != nil {
		switch err {
		case store.ErrExpertiseNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// updateExpertiseHandler godoc
//
//	@Summary		Updates a expertise field
//	@Description	Updates a expertise by ID
//	@Tags			expertise
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Expertise ID"
//	@Param			payload	body		UpdateExpertisePayload	true	"Expertise payload"
//	@Success		200		{object}	store.Expertise
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/expertise/{id} [patch]
func (app *application) updateExpertiseHandler(w http.ResponseWriter, r *http.Request) {
	expertiseID, err := strconv.ParseInt(chi.URLParam(r, "expertiseID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload UpdateExpertisePayload
	err = ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	expertise, err := app.store.Expertise.GetByID(r.Context(), expertiseID)
	if err != nil {
		switch err {
		case store.ErrExpertiseNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	
	if payload.Name != nil {
		expertise.Name = *payload.Name
	}
	if payload.Icon_svg != nil {
		expertise.Icon_svg = *payload.Icon_svg
	}

	err = app.store.Expertise.Update(r.Context(), expertise)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	err = JsonResponse(w, http.StatusOK, expertise)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getExpertiseHandler godoc
//
//	@Summary		Fetches all expertise
//	@Description	Fetches all expertise 
//	@Tags			expertise
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	store.Expertise
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/expertise [get]
func (app *application) getExpertiseHandler(w http.ResponseWriter, r *http.Request) {
	expertise, err := app.store.Expertise.Get(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusOK, expertise); err != nil {
		app.internalServerError(w, r, err)
	}
}