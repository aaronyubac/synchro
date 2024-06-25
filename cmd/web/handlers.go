package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"synchro/internal/models"
	"synchro/internal/validator"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	
	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	events, err := app.events.GetUserEvents(userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Events = events

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}


func (app *application) eventView(w http.ResponseWriter, r *http.Request) {

	eventId, err := strconv.Atoi(r.PathValue("event_id"))
	if err != nil || eventId < 1 {
		app.serverError(w, r, err)
		return 
	}

	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userId == 0 {
		http.Redirect(w, r, "user/login", http.StatusSeeOther)
		return
	}

	app.eventViewRenderer(w, r, userId, eventId, unavailabilityForm{})
	
	// event, err := app.events.GetEvent(userId, eventId)
	// if err != nil {
	// 	if errors.Is(err, models.ErrNoRecord) {
	// 		app.clientError(w, http.StatusNotFound) // CREATE app.NotFound as a wrapper around clientError?
	// 	} else {
	// 		app.serverError(w, r, err)
	// 	}
	// 	return
	// }

	// data := app.newTemplateData(r)
	// data.Event = event
	// data.Form = unavailabilityForm{}

	// app.render(w, r, http.StatusOK, "view.tmpl.html", data)

	
}

type unavailabilityForm struct {
	Date string
	AllDay string
	Start string
	End string
	validator.Validator
}

func (app *application) unavailabilityAdd(w http.ResponseWriter, r *http.Request) {

	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	eventId, err := strconv.Atoi(r.PathValue("event_id"))	
	if err != nil || eventId < 1 {
		app.serverError(w, r, err)
		return 
	}

	
	form := unavailabilityForm{}

	err = r.ParseForm();
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	dateStr := r.PostForm.Get("unavailability-date")
	startStr := r.PostForm.Get("unavailability-time-start")
	endStr := r.PostForm.Get("unavailability-time-end")
	unavailabilityAllDayStr := r.PostForm.Get("unavailability-all-day")


	form = unavailabilityForm{
		Date: dateStr,
		AllDay: unavailabilityAllDayStr,
		Start: startStr,
		End: endStr,
	}

	dateLayout := "2006-01-02 -0700"
	parsedDate, err := time.Parse(dateLayout, fmt.Sprintf("%s -0700", dateStr))
	if err != nil {
		form.AddFieldError("unavailabilityDate", "Select a valid date")
	}


	unavailabilityAllDayBool := false;

	if unavailabilityAllDayStr != "" {
		unavailabilityAllDayBool, err = strconv.ParseBool(unavailabilityAllDayStr)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}	
	}

	parsedStart := time.Time{}
	parsedEnd := time.Time{}

	if !unavailabilityAllDayBool {
	

		// MAKE TIMEZONE NOT HARDCODED
		timeLayout := "15:04 -0700 PDT"
		parsedStart, err = time.Parse(timeLayout, fmt.Sprintf("%s -0700 PDT", startStr))
		if err != nil {
			form.CheckField(validator.NotBlank(startStr), "time", "Enter a valid time")
		}

		parsedEnd, err = time.Parse(timeLayout, fmt.Sprintf("%s -0700 PDT", endStr))
		if err != nil {
			form.CheckField(validator.NotBlank(startStr), "time", "Enter a valid time")
		}

		
		if validator.NotBlank(startStr) && validator.NotBlank(endStr) {
		form.CheckField(validator.UnavailabilityTimeRange(parsedStart, parsedEnd), "time", "Enter a valid time range")
		}
		// check for user overlapping with previous unavailabilities
	}
	
		form.CheckField(validator.UnavailabilityNotPassed(parsedDate, parsedStart, parsedEnd), "unavailabilityDate", "Selected time/date has passed")

	if !form.Valid() {
		
		app.eventViewRenderer(w, r, userId, eventId, form)
		return
	}

	app.unavailabilities.Add(userId, eventId, parsedDate, form.Start, form.End, unavailabilityAllDayBool)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.eventViewRenderer(w, r, userId, eventId, unavailabilityForm{Date: dateStr})

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

	eventId, err := app.events.Create(form.Name, form.Details)
	if err != nil {
		app.serverError(w, r, err)
	}

	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userId == 0 {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	err = app.events.Join(userId, eventId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Event successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/event/%d", eventId), http.StatusSeeOther)
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

	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userId == 0 {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

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

	eventId, err := strconv.Atoi(form.EventID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.events.Join(userId, eventId)
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

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

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

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {

	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You have been logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)


}



