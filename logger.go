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
	Payload   map[string]interface{}
	Name      string
	TimeTaken time.Duration
}

var logger *logging.Logger
var dlpClient *dlp.Client
var ctx context.Context

func init() {
	ctx = context.Background()

	// Creates a client.
	client, err := logging.NewClient(ctx, "river-direction-210022")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	dlpClient, err = dlp.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Selects the log to write to.
	logger = client.Logger("blue-green-log")
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get time that request was received
		start := time.Now()

		// Serve the request
		inner.ServeHTTP(w, r)

		// Use DLP to mask sensative PII data
		r.ParseForm()
		formvalues, _ := json.Marshal(r.PostForm)

		buf := new(bytes.Buffer)
		var redactedForm map[string]interface{}
		if len(formvalues) > 0 {
			mask(buf, dlpClient, "river-direction-210022", string(formvalues), []string{}, "#", 0)
			json.Unmarshal(buf.Bytes(), &redactedForm)
		}

		// Log the request
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
