package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err)

	JSONError(w, http.StatusInternalServerError, "server encouter a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request response error", "method", r.Method, "path", r.URL.Path, "error", err)
	JSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err)

	JSONError(w, http.StatusNotFound, "resource not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("conflict error", "method", r.Method, "path", r.URL.Path, "error", err)

	JSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err)

	JSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedbasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	JSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) forbidden(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnf("user forbidden", "method", r.Method, "path", r.URL.Path)

	JSONError(w, http.StatusForbidden, "forbidden")
}
