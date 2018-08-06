// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command custommetric creates a custom metric and writes TimeSeries value
// to it. It writes a GAUGE measurement, which is a measure of value at a
// specific point in time. This means the startTime and endTime of the interval
// are the same. To make it easier to see the output, a random value is written.
// When reading the TimeSeries back, a window of the last 5 minutes is used.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/monitoring/v3"
)

const metricType = "custom.googleapis.com/deid_rate"

var dlp_rate int64

var (
	containerName string
	podName       string
	namespaceName string
)

func init() {
	dlp_rate = 0

	containerName = os.Getenv("HOSTNAME")
	podName = os.Getenv("POD_NAME")
	namespaceName = os.Getenv("NAMESPACE_NAME")
}

func projectResource(projectID string) string {
	return "projects/" + projectID
}

func metricCount() {
	dlp_rate = dlp_rate + 1
}

// createCustomMetric creates a custom metric specified by the metric type.
func createCustomMetric(s *monitoring.Service, projectID, metricType string) error {
	ld := monitoring.LabelDescriptor{Key: "my_label", ValueType: "STRING", Description: "Not in use"}
	md := monitoring.MetricDescriptor{
		Type:        metricType,
		Labels:      []*monitoring.LabelDescriptor{&ld},
		MetricKind:  "GAUGE",
		ValueType:   "INT64",
		Unit:        "items",
		Description: "Rate of de-identification",
		DisplayName: "De-identification Rate",
	}
	resp, err := s.Projects.MetricDescriptors.Create(projectResource(projectID), &md).Do()
	if err != nil {
		return fmt.Errorf("Could not create custom metric: %v", err)
	}

	_ = resp
	// log.Printf("createCustomMetric: %s\n", formatResource(resp))
	return nil
}

// getCustomMetric reads the custom metric created.
func getCustomMetric(s *monitoring.Service, projectID, metricType string) (*monitoring.ListMetricDescriptorsResponse, error) {
	resp, err := s.Projects.MetricDescriptors.List(projectResource(projectID)).
		Filter(fmt.Sprintf("metric.type=\"%s\"", metricType)).Do()
	if err != nil {
		return nil, fmt.Errorf("Could not get custom metric: %v", err)
	}

	_ = resp
	// log.Printf("getCustomMetric: %s\n", formatResource(resp))
	return resp, nil
}

// writeTimeSeriesValue writes a value for the custom metric created
func writeTimeSeriesValue(s *monitoring.Service, projectID, metricType string, metricValue int64) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	timeseries := monitoring.TimeSeries{
		Metric: &monitoring.Metric{
			Type: metricType,
			Labels: map[string]string{
				"my_label": "some_label",
			},
		},
		Resource: &monitoring.MonitoredResource{
			Labels: map[string]string{
				"container_name": containerName,
				"cluster_name":   "dlp-local-cluster",
				"location":       "us-central1",
				"namespace_name": namespaceName,
				"pod_name":       podName,
			},
			Type: "k8s_container",
		},
		Points: []*monitoring.Point{
			{
				Interval: &monitoring.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoring.TypedValue{
					Int64Value: &metricValue,
				},
			},
		},
	}

	createTimeseriesRequest := monitoring.CreateTimeSeriesRequest{
		TimeSeries: []*monitoring.TimeSeries{&timeseries},
	}

	// log.Printf("writeTimeseriesRequest: %s\n", formatResource(createTimeseriesRequest))
	_, err := s.Projects.TimeSeries.Create(projectResource(projectID), &createTimeseriesRequest).Do()
	if err != nil {
		return fmt.Errorf("Could not write time series value, %v ", err)
	}
	return nil
}

// readTimeSeriesValue reads the TimeSeries for the value specified by metric type in a time window from the last 5 minutes.
func readTimeSeriesValue(s *monitoring.Service, projectID, metricType string) error {
	startTime := time.Now().UTC().Add(time.Minute * -5)
	endTime := time.Now().UTC()
	resp, err := s.Projects.TimeSeries.List(projectResource(projectID)).
		Filter(fmt.Sprintf("metric.type=\"%s\"", metricType)).
		IntervalStartTime(startTime.Format(time.RFC3339Nano)).
		IntervalEndTime(endTime.Format(time.RFC3339Nano)).
		Do()
	if err != nil {
		return fmt.Errorf("Could not read time series value, %v ", err)
	}

	_ = resp
	// log.Printf("readTimeseriesValue: %s\n", formatResource(resp))
	return nil
}

func createService(ctx context.Context) (*monitoring.Service, error) {
	hc, err := google.DefaultClient(ctx, monitoring.MonitoringScope)
	if err != nil {
		return nil, err
	}
	s, err := monitoring.New(hc)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func pushMetrics(projectID string) {
	ctx := context.Background()
	s, err := createService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Create the metric.
	if err := createCustomMetric(s, projectID, metricType); err != nil {
		log.Fatal(err)
	}

	// Wait until the new metric can be read back.
	for {
		resp, err := getCustomMetric(s, projectID, metricType)
		if err != nil {
			log.Fatal(err)
		}
		if len(resp.MetricDescriptors) != 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}

	go func() {
		for range time.Tick(30 * time.Second) {
			// Write a TimeSeries value for that metric
			if err := writeTimeSeriesValue(s, projectID, metricType, dlp_rate); err != nil {
				log.Fatal(err)
			}
			dlp_rate = 0
		}
	}()

}

// formatResource marshals a response object as JSON.
func formatResource(resource interface{}) []byte {
	b, err := json.MarshalIndent(resource, "", "    ")
	if err != nil {
		panic(err)
	}
	return b
}
