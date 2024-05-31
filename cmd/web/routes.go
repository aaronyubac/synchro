package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticate)

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignupForm))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLoginForm))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /", protected.ThenFunc(app.home))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /event/create", protected.ThenFunc(app.eventCreateForm))
	mux.Handle("POST /event/create", protected.ThenFunc(app.eventCreatePost))
	mux.Handle("GET /event/{id}", protected.ThenFunc(app.eventView))
	mux.Handle("GET /event/join", protected.ThenFunc(app.eventJoinGet))
	mux.Handle("POST /event/join", protected.ThenFunc(app.eventJoinPost))

	return mux
}