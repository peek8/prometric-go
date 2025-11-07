/*
 * Copyright (c) 2025 peek8.io
 *
 * Created Date: Monday, November 3rd 2025, 3:23:33 pm
 * Author: Md. Asraful Haque
 *
 */

package prometrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HTTPMetricName string

const (
	HttpRequestsTotalMetric    HTTPMetricName = "http_requests_total"
	HttpRequestDurationMetric  HTTPMetricName = "http_request_duration_seconds"
	HttpRequestsInFlightMetric HTTPMetricName = "http_requests_in_flight"
	HttpRequestSizeMetric      HTTPMetricName = "http_request_size_bytes"
	HttpResponseSizeMetric     HTTPMetricName = "http_response_size_bytes"
)

var (
	// HttpRequestsTotal counts the total number of HTTP requests processed by the application.
	// It is labeled with the request path, HTTP method, and response status code.
	//
	// Typical usage:
	//
	//	HttpRequestsTotal.WithLabelValues("/api/v1/person", "GET", "200").Inc()
	//
	// Metric type: CounterVec
	HttpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: string(HttpRequestsTotalMetric),
		Help: "Total number of HTTP requests processed, labeled by status code and method.",
	}, []string{"path", "method", "code"})

	// HttpRequestDuration measures the duration of HTTP requests in seconds.
	// It is labeled by path, method, and status code, and uses the default Prometheus histogram buckets.
	//
	// Example usage:
	//
	//	timer := prometheus.NewTimer(HttpRequestDuration.WithLabelValues("/api/v1/person", "GET", "200"))
	//	defer timer.ObserveDuration()
	//
	// Metric type: HistogramVec
	HttpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    string(HttpRequestDurationMetric),
		Help:    "Histogram of HTTP request durations in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"path", "method", "code"})

	// HttpRequestsInFlight reports the number of HTTP requests currently being served.
	// It is labeled by request path.
	//
	// Typical usage:
	//
	//	HttpRequestsInFlight.WithLabelValues("/api/v1/person").Inc()
	//	defer HttpRequestsInFlight.WithLabelValues("/api/v1/person").Dec()
	//
	// Metric type: GaugeVec
	HttpRequestsInFlight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: string(HttpRequestsInFlightMetric),
		Help: "Number of HTTP requests currently being handled.",
	}, []string{"path"})

	// HttpRequestSize records the size of incoming HTTP requests in bytes.
	// It is labeled by path, method, and response code, and uses exponential buckets
	// starting at 100 bytes, growing by a factor of 10 up to 10^5.
	//
	// Example usage:
	//
	//	HttpRequestSize.WithLabelValues("/api/v1/person", "POST", "201").Observe(float64(req.ContentLength))
	//
	// Metric type: HistogramVec
	HttpRequestSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    string(HttpRequestSizeMetric),
		Help:    "Size of incoming HTTP requests in bytes.",
		Buckets: prometheus.ExponentialBuckets(100, 10, 5),
	}, []string{"path", "method", "code"})

	// HttpResponseSize records the size of outgoing HTTP responses in bytes.
	// It is labeled by path, method, and status code, and uses exponential buckets
	// starting at 100 bytes, growing by a factor of 10 up to 10^5.
	//
	// Example usage:
	//
	//	HttpResponseSize.WithLabelValues("/api/v1/person", "GET", "200").Observe(float64(respSize))
	//
	// Metric type: HistogramVec
	HttpResponseSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    string(HttpResponseSizeMetric),
		Help:    "Size of outgoing HTTP responses in bytes.",
		Buckets: prometheus.ExponentialBuckets(100, 10, 5),
	}, []string{"path", "method", "code"})
)
