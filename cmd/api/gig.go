package main

import (
	// "context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	chi "github.com/go-chi/chi/v5"
)

var ErrMissingExpertise = errors.New("expertise is not found")

type gigKey string

const gigCtx mentorKey = "gig"

type RegisterGigPayload struct {
	Title       string   `json:"title" validate:"required,max=100"`
	Description string   `json:"description" validate:"required"`
	Expertise   string   `json:"expertise" validate:"required"`
	Discipline  []string `json:"discipline" validate:"required"`
}

type UpdateGigPayload struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Expertise   *string   `json:"expertise"`
	Discipline  *[]string `json:"discipline"`
}

// createGigHandler godoc
//
//	@Summary		Creates a gig
//	@Description	Creates a gig
//	@Tags			gig
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterGigPayload	true	"Gig payload"
//	@Success		201		{object}	store.Gig
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/gigs/create [post]
func (app *application) createGigHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)

	var payload RegisterGigPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	gig := &store.Gig{
		Userid: user.ID,
		// Mentorid:    payload.MentorID,
		Title:       payload.Title,
		Description: payload.Description,
		Expertise:   payload.Expertise,
		Discipline:  payload.Discipline,
	}

	err = app.store.Gig.CreateGig(r.Context(), gig)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusCreated, gig)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getAllGigsHandler godoc
//
//	@Summary		Get all gigs
//	@Description	Get all gigs
//	@Tags			gig
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]store.Gig
//	@Failure		400	{object}	error
//	@Failure		401	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/gigs [get]
func (app *application) getAllGigsHandler(w http.ResponseWriter, r *http.Request) {
	gigs, err := app.store.Gig.GetAllGigs(r.Context(), 10, 0)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, gigs)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getGigsByExpertiseHandler godoc
//
//	@Summary		Get gigs by expertise
//	@Description	Get gigs by expertise
//	@Tags			gig
//	@Accept			json
//	@Produce		json
//	@Param			expertise	path		string	true	"Expertise"
//	@Success		200			{object}	[]store.Gig
//	@Failure		400			{object}	error
//	@Failure		401			{object}	error
//	@Failure		404			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/gigs/expertise/{expertise} [get]
func (app *application) getGigsByExpertiseHandler(w http.ResponseWriter, r *http.Request) {
	expertiseName := chi.URLParam(r, "expertise")
	if expertiseName == "" {
		app.badRequestResponse(w, r, ErrMissingExpertise)
		return
	}

	gigs, err := app.store.Gig.GetGigsByExpertise(r.Context(), expertiseName)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, gigs)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getGigByIDHandler godoc
//
//	@Summary		Get gig by ID
//	@Description	Get gig by ID
//	@Tags			gig
//	@Accept			json
//	@Produce		json
//	@Param			gigID	path		int	true	"Gig ID"
//	@Success		200		{object}	store.Gig
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/gigs/{gigID} [get]
func (app *application) getGigByIDHandler(w http.ResponseWriter, r *http.Request) {
	gigID, err := strconv.ParseInt(chi.URLParam(r, "gigID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// gigs := getGigFromCtx(r)

	gig, err := app.store.Gig.GetGigByID(r.Context(), gigID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, gig)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// updateGigHandler godoc
//
//	@Summary		Update gig
//	@Description	Update gig
//	@Tags			gig
//	@Accept			json
//	@Produce		json
//	@Param			gigID	path		int64				true	"Gig ID"
//	@Param			payload	body		UpdateGigPayload	true	"Gig payload"
//	@Success		200		{object}	store.Gig
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/gigs/{gigID} [patch]
func (app *application) updateGigHandler(w http.ResponseWriter, r *http.Request) {
	gigID, err := strconv.ParseInt(chi.URLParam(r, "gigID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	gig, err := app.store.Gig.GetGigByID(r.Context(), gigID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	var payload UpdateGigPayload
	err = ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Since all fields are optional, we only validate if at least one field is provided
	hasUpdates := payload.Title != nil || payload.Description != nil ||
		payload.Expertise != nil || payload.Discipline != nil
	if !hasUpdates {
		app.badRequestResponse(w, r, errors.New("no fields provided for update"))
		return
	}

	// Update only provided fields
	if payload.Title != nil {
		gig.Title = *payload.Title
	}
	if payload.Description != nil {
		gig.Description = *payload.Description
	}
	if payload.Expertise != nil {
		gig.Expertise = *payload.Expertise
	}
	if payload.Discipline != nil {
		gig.Discipline = *payload.Discipline
	}

	err = app.store.Gig.UpdateGig(r.Context(), gig)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, gig)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// deleteGigHandler godoc
//
//	@Summary		Delete gig
//	@Description	Delete gig
//	@Tags			gig
//	@Accept			json
//	@Produce		json
//	@Param			gigID	path		int	true	"Gig ID"
//	@Success		200		{object}	nil
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/gigs/{gigID} [delete]
func (app *application) deleteGigHandler(w http.ResponseWriter, r *http.Request) {
	gigID, err := strconv.ParseInt(chi.URLParam(r, "gigID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.Gig.DeleteGig(r.Context(), gigID)
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

// func (app *application) gigContextMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		idParam := chi.URLParam(r, "expertiseID")
// 		id, err := strconv.ParseInt(idParam, 10, 64)
// 		if err != nil {
// 			app.internalServerError(w, r, err)
// 			return
// 		}

// 		ctx := r.Context()

// 		gig, err := app.store.Gig.GetGigByID(ctx, id)
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, store.ErrNotFound):
// 				app.notFoundResponse(w, r, err)
// 			default:
// 				app.internalServerError(w, r, err)
// 			}
// 			return
// 		}

// 		ctx = context.WithValue(ctx, gigCtx, gig)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// func getGigFromCtx(r *http.Request) *store.Gig {
// 	gig, _ := r.Context().Value(gigCtx).(*store.Gig)
// 	return gig
// }
