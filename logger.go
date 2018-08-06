package main

import (
	"bytes"
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"time"

	dlp "cloud.google.com/go/dlp/apiv2"
	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

type MyEntry struct {
	Method    string
	Request   string
	Payload   interface{}
	Name      string
	TimeTaken int64
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

	pushMetrics(projectId)
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get time that request was received
		start := time.Now()

		// Serve the request
		inner.ServeHTTP(w, r)

		var redactedForm interface{}
		go func() {
			timeTaken := int64(time.Since(start) / time.Millisecond)

			if r.RequestURI == "/support" {
				// Parse data from form to string
				r.ParseForm()
				formEntries, _ := json.Marshal(r.PostForm)

				// Create byte buffer to hold redacted data
				buf := new(bytes.Buffer)

				if len(formEntries) > 0 {
					// Use DLP to mask sensative PII data
					mask(buf, dlpClient, projectId, string(formEntries), []string{}, "#", 0)
					//kms := "projects/river-direction-210022/locations/global/keyRings/dlp/cryptoKeys/dlp-demo"
					//deidentifyFPE(buf, dlpClient, projectId, string(formEntries), []string{}, "dek.enc", kms, "surrogateInfoType-test")
					json.Unmarshal(buf.Bytes(), &redactedForm)
				}

				fmt.Println("redactedForm: " + buf.String())
				fmt.Println("formEntries: " + string(formEntries))

				if buf.String() != string(formEntries) {
					metricCount()
				}
			}

			// Log the request along with redacted data
			if r.RequestURI != "/status" {
				logger.Log(logging.Entry{
					Severity: logging.Info,
					Resource: &monitoredres.MonitoredResource{
						Labels: map[string]string{
							"container_name": containerName,
							"cluster_name":   "dlp-local-cluster",
							"location":       "us-central1",
							"namespace_name": namespaceName,
							"pod_name":       podName,
						},
						Type: "k8s_container",
					},
					Payload: MyEntry{
						Method:    r.Method,
						Request:   r.RequestURI,
						Payload:   redactedForm,
						Name:      name,
						TimeTaken: timeTaken,
					},
				})
			}
		}()
	})
}
