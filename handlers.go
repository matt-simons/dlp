package main

import (
	"html/template"
	"net/http"
	"os"
)

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-store")

	// Load in form template
	tmpl := template.Must(template.ParseFiles("index.html"))

	if len(os.Getenv("AWS")) > 0 {
		tmpl.Execute(w, struct{ Aws bool }{true})
		return
	}

	tmpl.Execute(w, nil)
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
