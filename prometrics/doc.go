// Package prometrics provides simple, composable Prometheus instrumentation
// for HTTP servers and application-level metrics.
//
// # Overview
//
// This package offers ready-to-use middleware wrappers and helper utilities
// to expose standardized Prometheus metrics for HTTP handlers, application
// health, and business logic metrics (such as object counts or operation latency).
//
// It is designed to reduce boilerplate while keeping flexibility for advanced
// users. The library wraps Prometheus client primitives such as CounterVec,
// GaugeVec, and HistogramVec, and integrates seamlessly with net/http or Gin.
//
// Example usage:
//
//	import (
//	    "net/http"
//	    "github.com/peek8/prometric-go/prometrics"
//	)
//
//	func main() {
//	    mux := http.NewServeMux()
//	    mux.Handle("/person", prometrics.InstrumentHttpHandler("person_handler", http.HandlerFunc(PersonHandler)))
//
//	    // Expose metrics endpoint
//	    http.Handle("/metrics", promhttp.Handler())
//
//	    http.ListenAndServe(":8080", mux)
//	}
//
//	func PersonHandler(w http.ResponseWriter, r *http.Request) {
//	    // your CRUD logic
//	}
//
// # Provided Metrics
//
// By default, all of the following HTTP metrics are exposed:
//   - http_requests_total{path,method,code}
//   - http_request_duration_seconds{path,method,code}
//   - http_in_flight_requests{path}
//   - http_request_size_bytes{path,method,code}
//   - http_response_size_bytes{path,method,code}
//
// Additionally, application-level gauges or counters can be created dynamically
// using the MetricFactory API for business metrics.
package prometrics
