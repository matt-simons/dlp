package main

import (
	"html/template"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type ContactDetails struct {
	Email   string
	Subject string
	Message string
}

func MakePayment(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("forms.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	details := ContactDetails{
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}

	// do something with details
	_ = details

	tmpl.Execute(w, struct{ Success bool }{true})
}

func ChallengeBan(w http.ResponseWriter, r *http.Request) {

}

func Phishing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"path":   path,
		"method": r.Method,
		"body":   r.PostForm,
	}).Info("New phishing attempt made")

}
