/*
 * Copyright (c) 2025 peek8.io
 *
 * Created Date: Monday, November 3rd 2025, 3:23:33 pm
 * Author: Md. Asraful Haque
 *
 */

// Package prometrics provide middleware to expose common prometheus metrics for http application.
package prometrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type CurriableVec interface {
	prometheus.Collector
	MustCurryWith(labels prometheus.Labels) prometheus.Collector
}

type HTTPMetricName string

const (
	HttpRequestsTotalMetric    HTTPMetricName = "http_requests_total"
	HttpRequestDurationMetric  HTTPMetricName = "http_request_duration_seconds"
	HttpRequestsInFlightMetric HTTPMetricName = "http_requests_in_flight"
	HttpRequestSizeMetric      HTTPMetricName = "http_request_size_bytes"
	HttpResponseSizeMetric     HTTPMetricName = "http_response_size_bytes"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: string(HttpRequestsTotalMetric),
		Help: "Total number of HTTP requests processed, labeled by status code and method.",
	}, []string{"handler", "method", "code"})

	HttpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    string(HttpRequestDurationMetric),
		Help:    "Histogram of HTTP request durations in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	HttpRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: string(HttpRequestsInFlightMetric),
		Help: "Number of HTTP requests currently being handled.",
	})

	HttpRequestSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    string(HttpRequestSizeMetric),
		Help:    "Size of incoming HTTP requests in bytes.",
		Buckets: prometheus.ExponentialBuckets(100, 10, 5),
	}, []string{"handler", "method", "code"})

	HttpResponseSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    string(HttpResponseSizeMetric),
		Help:    "Size of outgoing HTTP responses in bytes.",
		Buckets: prometheus.ExponentialBuckets(100, 10, 5),
	}, []string{"handler", "method", "code"})
)
