package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"synchro/internal/models"
)
func (app *application) clientError(w http.ResponseWriter, status int) {
	
	http.Error(w, http.StatusText(status), status)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {

	var(
		method = r.Method
		uri = r.URL.RequestURI
		trace = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}


func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}


func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData {
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {

	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) eventViewRenderer(w http.ResponseWriter, r *http.Request, userId, eventId int, form any) {
	event, err := app.events.GetEvent(userId, eventId)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.clientError(w, http.StatusNotFound) // CREATE app.NotFound as a wrapper around clientError?
			} else {
				app.serverError(w, r, err)
			}
		return
		}


		data := app.newTemplateData(r)
		data.Event = event
		data.Form = form

		app.render(w, r, http.StatusBadRequest, "view.tmpl.html", data)
	}