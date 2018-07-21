package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/cloudflare/cloudflare-go"
	"github.com/gorilla/mux"
)

func MakePayment(w http.ResponseWriter, r *http.Request) {

}

func ChallengeBan(w http.ResponseWriter, r *http.Request) {
	ip := net.ParseIP(r.Header.Get("CF-Connecting-IP"))

	rule := cloudflare.AccessRule{
		Notes: "Auto ban by the bouncer",
		Mode:  "challenge",
		Configuration: cloudflare.AccessRuleConfiguration{
			Target: "ip",
			Value:  ip.String(),
		},
	}

	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"rule": rule,
	}).Info("Creating new challenge for IP")

	res, err := api.CreateZoneAccessRule(zone.ID, rule)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", res)
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
