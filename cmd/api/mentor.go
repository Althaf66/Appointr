package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	"github.com/go-chi/chi/v5"
)

type mentorKey string

const mentorCtx mentorKey = "mentor"

type RegisterMentorPayload struct {
	Name     string   `json:"name" validate:"required,max=40"`
	Country  string   `json:"country" validate:"required"`
	Language []string `json:"language" validate:"required"`
}

// createMentorHandler godoc
//
//	@Summary		Creates a mentor
//	@Description	Creates a mentor
//	@Tags			mentor
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterMentorPayload	true	"Mentor payload"
//	@Success		201		{object}	store.Mentor
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/mentors/create [post]
func (app *application) createMentorHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)
	var payload RegisterMentorPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	mentor := &store.Mentor{
		Userid:   user.ID,
		Name:     payload.Name,
		Country:  payload.Country,
		Language: payload.Language,
	}

	err = app.store.Mentor.CreateMentor(r.Context(), mentor)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusCreated, mentor)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getMentorsHandler godoc
//
//	@Summary		Get all mentors
//	@Description	Get all mentors
//	@Tags			mentor
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int	false	"Limit"
//	@Param			offset	query		int	false	"Offset"
//	@Success		200		{object}	[]store.Mentor
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/mentors [get]
func (app *application) getMentorsHandler(w http.ResponseWriter, r *http.Request) {
	mentors, err := app.store.Mentor.GetAllMentors(r.Context(), 10, 0)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, mentors)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getMentorByNameHandler godoc
//
//	@Summary		Get mentor by name
//	@Description	Get mentor by name
//	@Tags			mentor
//	@Accept			json
//	@Produce		json
//	@Param			mentorName	path		string	true	"Mentor Name"
//	@Success		200			{object}	[]store.Mentor
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/mentors/name/{mentorName} [get]
func (app *application) getMentorByNameHandler(w http.ResponseWriter, r *http.Request) {
	mentorName := chi.URLParam(r, "mentorName")

	mentors, err := app.store.Mentor.GetMentorByName(r.Context(), mentorName)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, mentors)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getMentorByExpertiseHandler godoc
//
//	@Summary		Get mentor by expertise
//	@Description	Get mentor by expertise
//	@Tags			mentor
//	@Accept			json
//	@Produce		json
//	@Param			expertise	path		string	true	"Expertise"
//	@Success		200			{object}	[]store.Mentor
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/mentors/exp/{expertise} [get]
func (app *application) getMentorByExpertiseHandler(w http.ResponseWriter, r *http.Request) {
	expertise := chi.URLParam(r, "expertise")

	mentors, err := app.store.Mentor.GetMentorsByExpertise(r.Context(), expertise)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, mentors)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getMentorByDisciplineHandler godoc
//
//	@Summary		Get mentor by discipline
//	@Description	Get mentor by discipline
//	@Tags			mentor
//	@Accept			json
//	@Produce		json
//	@Param			discipline	path		string	true	"Expertise"
//	@Success		200			{object}	[]store.Mentor
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/mentors/dis/{discipline} [get]
func (app *application) getMentorByDisciplineHandler(w http.ResponseWriter, r *http.Request) {
	discipline := chi.URLParam(r, "discipline")

	mentors, err := app.store.Mentor.GetMentorsByDiscipline(r.Context(), discipline)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, mentors)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getMentorByIDHandler godoc
//
//	@Summary		Get mentor by ID
//	@Description	Get mentor by ID
//	@Tags			mentor
//	@Accept			json
//	@Produce		json
//	@Param			mentorID	path		int	true	"Mentor ID"
//	@Success		200			{object}	store.Mentor
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/mentors/{mentorID} [get]
func (app *application) getMentorByIDHandler(w http.ResponseWriter, r *http.Request) {
	mentor := getMentorFromCtx(r)

	mentors, err := app.store.Mentor.GetMentorByID(r.Context(), mentor.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, mentors)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getMentorByUserIDHandler godoc
//
//	@Summary		Get mentor by userID
//	@Description	Get mentor by userID
//	@Tags			mentor
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"Mentor ID"
//	@Success		200		{object}	store.Mentor
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/mentors/u/{userID} [get]
func (app *application) getMentorByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	// user := getUserfromCtx(r)
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	mentors, err := app.store.Mentor.GetMentorByUserID(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, mentors)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type UpdateMentorPayload struct {
	Name     *string   `json:"name" validate:"omitempty,max=40"`
	Country  *string   `json:"country" validate:"omitempty"`
	Language *[]string `json:"language" validate:"omitempty"`
}

// updateMentorHandler godoc
//
//	@Summary		Update mentor
//	@Description	Update mentor
//	@Tags			mentor
//	@Accept			json
//	@Produce		json
//	@Param			mentorID	path		int					true	"Mentor ID"
//	@Param			payload		body		UpdateMentorPayload	true	"Mentor payload"
//	@Success		200			{object}	store.Mentor
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/mentors/{mentorID} [patch]
func (app *application) updateMentorHandler(w http.ResponseWriter, r *http.Request) {
	mentor := getMentorFromCtx(r)

	var payload UpdateMentorPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	hasUpdates := payload.Name != nil || payload.Country != nil ||
		payload.Language != nil
	if !hasUpdates {
		app.badRequestResponse(w, r, errors.New("no fields provided for update"))
		return
	}

	if payload.Name != nil {
		mentor.Name = *payload.Name
	}
	if payload.Country != nil {
		mentor.Country = *payload.Country
	}
	if payload.Language != nil {
		mentor.Language = *payload.Language
	}

	err = app.store.Mentor.UpdateMentor(r.Context(), mentor)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusOK, mentor); err != nil {
		app.internalServerError(w, r, err)
	}
}

// deleteMentorHandler godoc
//
//	@Summary		Delete mentor
//	@Description	Delete mentor
//	@Tags			mentor
//	@Accept			json
//	@Produce		json
//	@Param			mentorID	path		int	true	"Mentor ID"
//	@Success		200			{object}	nil
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/mentors/{mentorID} [delete]
func (app *application) deleteMentorHandler(w http.ResponseWriter, r *http.Request) {
	mentor := getMentorFromCtx(r)

	err := app.store.Mentor.DeleteMentor(r.Context(), mentor.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, nil)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) mentorContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "mentorID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		mentor, err := app.store.Mentor.GetMentorByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, mentorCtx, mentor)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getMentorFromCtx(r *http.Request) *store.Mentor {
	mentor, _ := r.Context().Value(mentorCtx).(*store.Mentor)
	return mentor
}
