package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticate)

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignupForm))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLoginForm))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /", protected.ThenFunc(app.home))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("POST /event/create", protected.ThenFunc(app.eventCreatePost))
	mux.Handle("GET /event/{event_id}", protected.ThenFunc(app.eventView))
	mux.Handle("POST /event/{event_id}", protected.ThenFunc(app.unavailabilityAdd))
	mux.Handle("POST /unavailability/remove", protected.ThenFunc(app.unavailabilityRemove))
	mux.Handle("POST /event/join", protected.ThenFunc(app.eventJoinPost))

	return mux
}