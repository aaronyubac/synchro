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
	if userId == 0 {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	events, err := app.events.GetUserEvents(userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = eventCreateForm{}
	data.Events = events

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}


func (app *application) eventView(w http.ResponseWriter, r *http.Request) {

	// eventId, err := strconv.Atoi(r.PathValue("event_id"))
	// if err != nil || eventId < 1 {
	// 	app.serverError(w, r, err)
	// 	return 
	// }

	eventId := r.PathValue("event_id")

	userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userId == 0 {
		http.Redirect(w, r, "user/login", http.StatusSeeOther)
		return
	}

	app.eventViewRenderer(w, r, userId, eventId, unavailabilityForm{}, http.StatusOK)
	
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
	eventId := r.PathValue("event_id")

	
	form := unavailabilityForm{}

	err := r.ParseForm();
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

	dateLayout := "2006-01-02"
	parsedDate, err := time.Parse(dateLayout, dateStr)
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

	startDateTime := parsedDate
	endDateTime := parsedDate.Add(time.Hour * 23 + time.Minute * 59)

	if !unavailabilityAllDayBool {
	
		timeLayout := "15:04 -0700 PDT"
		parsedStartTime, err := time.Parse(timeLayout, fmt.Sprintf("%s -0700 PDT", startStr))
		if err != nil {
			form.CheckField(validator.NotBlank(startStr), "time", "Enter a valid time")
		}

		parsedEndTime, err := time.Parse(timeLayout, fmt.Sprintf("%s -0700 PDT", endStr))
		if err != nil {
			form.CheckField(validator.NotBlank(startStr), "time", "Enter a valid time")
		}

		// convert pdt to utc
		parsedStartTime = parsedStartTime.UTC()
		parsedEndTime = parsedEndTime.UTC()

		// check for user overlapping with previous unavailabilities
		
		if validator.NotBlank(startStr) && validator.NotBlank(endStr) {
		form.CheckField(validator.UnavailabilityTimeRange(parsedStartTime, parsedEndTime), "time", "Enter a valid time range")
		}

		form.CheckField(validator.PermittedMinutes(parsedStartTime), "time", "Minutes must be a multiple of 15")
		form.CheckField(validator.PermittedMinutes(parsedEndTime), "time", "Minutes must be a multiple of 15")

		startDateTime = parsedDate.Add(((time.Hour * 24) * time.Duration(parsedStartTime.Day() - 1))+ time.Hour * time.Duration(parsedStartTime.Hour()) + time.Minute * time.Duration(parsedStartTime.Minute()))
		endDateTime = parsedDate.Add(((time.Hour * 24) * time.Duration(parsedStartTime.Day() - 1)) + time.Hour * time.Duration(parsedEndTime.Hour()) + time.Minute * time.Duration(parsedEndTime.Minute()))
	}

	form.CheckField(validator.TimeNotPassed(startDateTime), "unavailabilityDate", "Selected time/date has passed")


	if !form.Valid() {
		form.AddNonFieldError("Failed to add unavailability")
		app.eventViewRenderer(w, r, userId, eventId, form, http.StatusUnprocessableEntity)
		return
	}

	app.unavailabilities.Add(userId, eventId, startDateTime.Format("2006-01-02 15:04:05"), endDateTime.Format("2006-01-2 15:04:05"), unavailabilityAllDayBool)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.eventViewRenderer(w, r, userId, eventId, unavailabilityForm{Date: dateStr}, http.StatusOK)

}

func (app *application) unavailabilityRemove(w http.ResponseWriter, r *http.Request) {
	
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

	formUserId, err := strconv.Atoi(r.PostForm.Get("userId"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if userId == formUserId {

		unavailabilityId, err := strconv.Atoi(r.PostForm.Get("unavailabilityId"))
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		app.unavailabilities.RemoveUserUnavailability(userId, unavailabilityId)

	} 

	eventId:= r.PostForm.Get("eventId")

	http.Redirect(w, r, fmt.Sprintf("/event/%s", eventId), http.StatusSeeOther)


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
		http.Redirect(w, r, "/", http.StatusBadRequest)
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

	http.Redirect(w, r, fmt.Sprintf("/event/%s", eventId), http.StatusSeeOther)
}

type eventJoinForm struct {
	EventID string
	validator.Validator
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
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	err = app.events.Join(userId, form.EventID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			form.AddNonFieldError("Invalid Event ID")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			
		} else if errors.Is(err, models.ErrDuplicateEvent) {
			form.AddNonFieldError("Already part of event")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			app.serverError(w, r, err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
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



