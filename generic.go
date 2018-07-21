package main

import (
	"html/template"
	"net/http"
)

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

func ContactSupport(w http.ResponseWriter, r *http.Request) {

	// Load in form template
	tmpl := template.Must(template.ParseFiles("forms.html"))

	// Show form
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	// Get details from submitted form
	details := ContactDetails{
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}

	// Do something with details
	_ = details

	// Show success page
	tmpl.Execute(w, struct{ Success bool }{true})
}

