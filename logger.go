package main

import (
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
)

var logger logging.Client
var ctx Context

func init() {
	ctx = context.Background()

	// Creates a client.
	client, err := logging.NewClient(ctx, "river-direction-210022")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Selects the log to write to.
	logger = client.Logger("blue-green-log")
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		if r.RequestURI != "/status" {
			logger.Log(logging.Entry{
				"%s\t%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				name,
				time.Since(start),
			})
		}
	})
}
