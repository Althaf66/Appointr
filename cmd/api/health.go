package main

import "net/http"

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"env":     app.config.env,
		"status":  "healthy",
		"version": "0.0.1",
	}
	err := JsonResponse(w, http.StatusOK, data)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
