package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Althaf66/Appointr/internal/store"
	"github.com/go-chi/chi/v5"
)

type RegisterSocialMediaPayload struct {
	Name string `json:"name" validate:"required"`
	Link string `json:"link"`
}

type UpdateSocialMediaPayload struct {
	Name *string `json:"name"`
	Link *string `json:"link"`
}

// createSocialMediaHandler godoc
//
//	@Summary		Create a new socialmedia
//	@Description	Create a new socialmedia
//	@Tags			socialmedia
//	@Accept			json
//	@Produce		json
//	@Param			socialmedia	body		RegisterSocialMediaPayload	true	"SocialMedia"
//	@Success		200			{object}	error
//	@Failure		400			{object}	error
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/socialmedia/create [post]
func (app *application) createSocialMediaHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)

	var payload RegisterSocialMediaPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	socialmedia := &store.SocialMedia{
		Userid: user.ID,
		Name:   payload.Name,
		Link:   payload.Link,
	}

	err = app.store.SocialMedia.CreateSocialMedia(r.Context(), socialmedia)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	err = JsonResponse(w, http.StatusCreated, socialmedia)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getSocialMediaByUserIDHandler godoc
//
//	@Summary		Get all socialmedia by user ID
//	@Description	Get all socialmedia by user ID
//	@Tags			socialmedia
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{array}		store.SocialMedia
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/socialmedia/u/{id} [get]
func (app *application) getSocialMediaByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	socialmedia, err := app.store.SocialMedia.GetSocialMediaByUserId(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, socialmedia)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getSocialMediaByIDHandler godoc
//
//	@Summary		Get socialmedia by ID
//	@Description	Get socialmedia by ID
//	@Tags			socialmedia
//	@Accept			json
//	@Produce		json
//	@Param			socialMediaID	path		int	true	"SocialMedia ID"
//	@Success		200				{object}	store.SocialMedia
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/socialmedia/{socialMediaID} [get]
func (app *application) getSocialMediaByIDHandler(w http.ResponseWriter, r *http.Request) {
	socialmediaID, err := strconv.ParseInt(chi.URLParam(r, "socialMediaID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	socialmedia, err := app.store.SocialMedia.GetSocialMediaById(r.Context(), socialmediaID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = JsonResponse(w, http.StatusOK, socialmedia)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// updateSocialMediaHandler godoc
//
// @Summary		Update socialmedia
// @Description	Update socialmedia
// @Tags			socialmedia
// @Accept			json
// @Produce		json
// @Param			socialMediaID	path		int							true	"SocialMedia ID"
// @Param			socialmedia		body		UpdateSocialMediaPayload	true	"SocialMedia"
// @Success		200				{object}	store.SocialMedia
// @Failure		404				{object}	error
// @Failure		500				{object}	error
// @Security		ApiKeyAuth
// @Router			/socialmedia/{socialMediaID} [patch]
func (app *application) updateSocialMediaHandler(w http.ResponseWriter, r *http.Request) {
	socialmediaID, err := strconv.ParseInt(chi.URLParam(r, "socialMediaID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	socialmedia, err := app.store.SocialMedia.GetSocialMediaById(r.Context(), socialmediaID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	var payload UpdateSocialMediaPayload
	err = ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Since all fields are optional, we only validate if at least one field is provided
	hasUpdates := payload.Name != nil || payload.Link != nil
	if !hasUpdates {
		app.badRequestResponse(w, r, errors.New("no fields provided for update"))
		return
	}

	// Update only provided fields
	if payload.Name != nil {
		socialmedia.Name = *payload.Name
	}
	if payload.Link != nil {
		socialmedia.Link = *payload.Link
	}

	err = app.store.SocialMedia.UpdateSocialMedia(r.Context(), socialmedia)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = JsonResponse(w, http.StatusOK, socialmedia)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// deleteGigHandler godoc
//
//	@Summary		Delete socialmedia
//	@Description	Delete socialmedia by id
//	@Tags			socialmedia
//	@Accept			json
//	@Produce		json
//	@Param			socialMediaID	path		int	true	"SocialMedia ID"
//	@Success		200				{object}	nil
//	@Failure		400				{object}	error
//	@Failure		401				{object}	error
//	@Failure		500				{object}	error
//	@Security		ApiKeyAuth
//	@Router			/socialmedia/{socialMediaID} [delete]
func (app *application) deleteSocialMediaHandler(w http.ResponseWriter, r *http.Request) {
	socialmediaID, err := strconv.ParseInt(chi.URLParam(r, "socialMediaID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.SocialMedia.DeleteSocialMedia(r.Context(), socialmediaID)
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
