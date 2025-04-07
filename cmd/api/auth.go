package main

import (
	// "crypto/rand"
	"crypto/sha256"
	// "encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/Althaf66/Appointr/internal/mailer"
	"github.com/Althaf66/Appointr/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	err := ReadJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = Validate.Struct(payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashtoken := hex.EncodeToString(hash[:])

	err = app.store.Users.CreateAndInvite(r.Context(), user, hashtoken, app.config.mail.exp)
	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}
	isProdenv := app.config.env == "production"

	status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdenv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)
		// rollback if any failure happens
		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}
		app.internalServerError(w, r, err)
		return
	}
	app.logger.Infow("Email sent", "status code", status)

	err = JsonResponse(w, http.StatusCreated, userWithToken)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// createTokenHandler godoc
//
//	@Summary		Creates a token
//	@Description	Creates a token for a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string					"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserTokenPayload
	if err := ReadJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrUserNotFound:
			app.unauthorizedErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := JsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}
}

// func generateState() string {
// 	b := make([]byte, 16)
// 	_, _ = rand.Read(b)
// 	return base64.URLEncoding.EncodeToString(b)
// }

// func (app *application) googleLogin(w http.ResponseWriter, r *http.Request) {
// 	state := generateState()
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "oauthstate",
// 		Value:    state,
// 		Expires:  time.Now().Add(1 * time.Hour),
// 		HttpOnly: true,
// 	})
// 	url := googleOauthConfig.AuthCodeURL(state)
// 	http.Redirect(w, r, url, http.StatusFound)
// }

// func (app *application) googleCallback(w http.ResponseWriter, r *http.Request) {
// 	var payload RegisterUserPayload
// 	// Verify state
// 	cookie, err := r.Cookie("oauthstate")
// 	if err != nil || r.URL.Query().Get("state") != cookie.Value {
// 		http.Error(w, "Invalid OAuth state", http.StatusUnauthorized)
// 		return
// 	}

// 	code := r.URL.Query().Get("code")
// 	token, err := googleOauthConfig.Exchange(r.Context(), code)
// 	if err != nil {
// 		http.Error(w, "Code exchange failed", http.StatusUnauthorized)
// 		return
// 	}

// 	client := googleOauthConfig.Client(r.Context(), token)
// 	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
// 	if err != nil {
// 		http.Error(w, "Failed to get user info", http.StatusUnauthorized)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	err = ReadJSON(w, r, &payload)
// 	if err != nil {
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	// Insert or update user in the database
// 	user, err := app.store.Users.GauthCreate(r.Context(), payload.Email, payload.Username)
// 	if err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	claims := jwt.MapClaims{
// 		"sub": user.ID,
// 		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
// 		"iat": time.Now().Unix(),
// 		"nbf": time.Now().Unix(),
// 		"iss": app.config.auth.token.iss,
// 		"aud": app.config.auth.token.iss,
// 	}

// 	tokens, err := app.authenticator.GenerateToken(claims)
// 	if err != nil {
// 		app.internalServerError(w, r, err)
// 		return
// 	}

// 	if err := JsonResponse(w, http.StatusCreated, tokens); err != nil {
// 		app.internalServerError(w, r, err)
// 	}
// }
