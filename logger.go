package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
)

type MyEntry struct {
	Method    string
	Request   string
	Payload   interface{}
	Name      string
	TimeTaken time.Duration
}

var projectId = "river-direction-210022"
var logger *logging.Logger
var dlpClient *dlp.Client

func init() {
	ctx := context.Background()

	// Creates a client.
	client, err := logging.NewClient(ctx, projectId)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Selects the log to write to.
	logger = client.Logger("blue-green-log")

	// Create Data Loss Protection client
	dlpClient, err = dlp.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get time that request was received
		start := time.Now()

		// Serve the request
		inner.ServeHTTP(w, r)

		// Parse data from form to string
		r.ParseForm()
		formEntries, _ := json.Marshal(r.PostForm)

		// Create byte buffer to hold redacted data 
		buf := new(bytes.Buffer)
		var redactedForm interface{}

		if len(formEntries) > 0 {
			// Use DLP to mask sensative PII data
			mask(buf, dlpClient, projectId, string(formEntries), []string{}, "#", 0)
			json.Unmarshal(buf.Bytes(), &redactedForm)
		}

		// Log the request along with redacted data
		if r.RequestURI != "/status" {
			logger.Log(logging.Entry{
				Severity: logging.Info,
				Payload: MyEntry{
					Method:    r.Method,
					Request:   r.RequestURI,
					Payload:   redactedForm,
					Name:      name,
					TimeTaken: time.Since(start),
				},
			})
		}
	})
}
