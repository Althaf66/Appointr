package main

import (
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	chi "github.com/go-chi/chi/v5"
)

type RegisterDisciplinePayload struct {
	Field    string `json:"field" `
	Subfield string `json:"subfield"`
}

type UpdateDisciplinePayload struct {
	Field    *string `json:"field" `
	Subfield *string `json:"subfield"`
}

// createDisciplineHandler godoc
//
//	@Summary		Creates a new discipline
//	@Description	creates a new discipline field
//	@Tags			discipline
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterDisciplinePayload	true	"discipline"
//	@Success		201		{object}	int64						"Discipline registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/discipline/create [post]
func (app *application) createDisciplineHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterDisciplinePayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	discipline := &store.Discipline{
		Field:    payload.Field,
		SubField: payload.Subfield,
	}

	err = app.store.Discipline.Create(r.Context(), discipline)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	disciplineID := discipline.ID
	err = JsonResponse(w, http.StatusCreated, disciplineID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

//	 getDisciplineHandlerByField godoc
//
//	@Summary		Fetches discipline
//	@Description	Fetches discipline by Field
//	@Tags			discipline
//	@Accept			json
//	@Produce		json
//	@Param			string	path		string	true	"Discipline Field"
//	@Success		200		{object}	store.Discipline
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/discipline/{string} [get]
func (app *application) getDisciplineHandlerByField(w http.ResponseWriter, r *http.Request) {
	disciplinefield := chi.URLParam(r, "disciplineField")

	discipline, err := app.store.Discipline.GetByField(r.Context(), disciplinefield)
	if err != nil {
		switch err {
		case store.ErrDisciplineNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := JsonResponse(w, http.StatusOK, discipline); err != nil {
		app.internalServerError(w, r, err)
	}
}

// deleteDisciplineHandler godoc
//
//	@Summary		Deletes a discipline field
//	@Description	Delete a discipline field by ID
//	@Tags			discipline
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Discipline ID"
//	@Success		204	{object}	string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/discipline/{id} [delete]
func (app *application) deleteDisciplineHandler(w http.ResponseWriter, r *http.Request) {
	disciplineID, err := strconv.ParseInt(chi.URLParam(r, "disciplineID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.Discipline.Delete(r.Context(), disciplineID)
	if err != nil {
		switch err {
		case store.ErrDisciplineNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// updateDisciplineHandler godoc
//
//	@Summary		Updates a discipline field
//	@Description	Updates a discipline by ID
//	@Tags			discipline
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Discipline ID"
//	@Param			payload	body		UpdateDisciplinePayload	true	"Discipline payload"
//	@Success		200		{object}	store.Discipline
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/discipline/{id} [patch]
func (app *application) updateDisciplineHandler(w http.ResponseWriter, r *http.Request) {
	disciplineID, err := strconv.ParseInt(chi.URLParam(r, "disciplineID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload UpdateDisciplinePayload
	err = ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	discipline, err := app.store.Discipline.GetByID(r.Context(), disciplineID)
	if err != nil {
		switch err {
		case store.ErrDisciplineNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if payload.Field != nil {
		discipline.Field = *payload.Field
	}
	if payload.Subfield != nil {
		discipline.SubField = *payload.Subfield
	}

	err = app.store.Discipline.Update(r.Context(), discipline)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	err = JsonResponse(w, http.StatusOK, discipline)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getDisciplineHandler godoc
//
//	@Summary		Fetches all discipline
//	@Description	Fetches all discipline
//	@Tags			discipline
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	store.Discipline
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/discipline [get]
func (app *application) getDisciplineHandler(w http.ResponseWriter, r *http.Request) {
	discipline, err := app.store.Discipline.Get(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusOK, discipline); err != nil {
		app.internalServerError(w, r, err)
	}
}
