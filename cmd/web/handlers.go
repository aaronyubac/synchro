package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"synchro/internal/models"
	"synchro/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}


func (app *application) eventView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.serverError(w, r, err)
		return 
	}
	
	event, err := app.events.Get(id)
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

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)

	
}

func (app *application) eventCreateForm(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Form = eventCreateForm{}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)

}

type eventCreateForm struct {
	Name string
	Details string
	validator.Validator
}

func (app *application) eventCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	name := r.PostForm.Get("name")
	details := r.PostForm.Get("details")

	form := eventCreateForm{
		Name: name,
		Details: details,
	}

		form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
		form.CheckField(validator.MaxChars(form.Name, 100), "name", "This field cannot be more than 100 characters long")
		form.CheckField(validator.MaxChars(form.Details, 1023), "details", "This field cannot be more than 1023 characters long")
	
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusBadRequest, "create.tmpl.html", data)
		return
	}

	id, err := app.events.Create(form.Name, form.Details)
	if err != nil {
		app.serverError(w, r, err)
	}

	app.sessionManager.Put(r.Context(), "flash", "Event successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/event/%d", id), http.StatusSeeOther)
}



func (app *application) userSignupForm(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Form = userSignupForm{}

	app.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

type userSignupForm struct {
	Name string
	Email string
	Password string
	validator.Validator
}


func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}


	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	form := userSignupForm{
		Name: name,
		Email: email,
		Password: password,
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be left blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be left blank")
	// Check if email is right format
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be left blank")
	// Minimum password length

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusBadRequest, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (app *application) userLoginForm(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	app.render(w, r, http.StatusOK, "login.tmpl.html", data)
}

type userLoginForm struct {
	Email string
	Password string
	validator.Validator

}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	form := userLoginForm{
		Email: email,
		Password: password,
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be left blank")
	// validator to check if in email format
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be left blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}
	

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
}

type eventJoinForm struct {
	EventID string
	validator.Validator
}

func (app *application) eventJoinGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = eventJoinForm{}
	app.render(w, r, http.StatusOK, "join.tmpl.html", data)
}

func (app *application) eventJoinPost(w http.ResponseWriter, r *http.Request) {

	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userID == 0 {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}

	r.ParseForm() // err???

	form := eventJoinForm{
		EventID: r.PostForm.Get("eventID"),
	}

	form.CheckField(validator.NotBlank(form.EventID), "eventID", "This field cannot be left blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "join.tmpl.html", data)
		return
	}

	eventID, err := strconv.Atoi(form.EventID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.events.Join(userID, eventID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {

			form.AddNonFieldError("Invalid Event ID")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusNotFound, "join.tmpl.html", data)
			
		} else if errors.Is(err, models.ErrDuplicateEvent) {

			form.AddNonFieldError("Already part of event")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "join.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	
	http.Redirect(w, r, fmt.Sprintf("/event/%s", form.EventID), http.StatusSeeOther)
}

