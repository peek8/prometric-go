package prometrics_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/peek8/prometric-go/prometrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ExampleInstrumentHttpHandler demonstrates how to instrument a standard net/http handler
// with Prometheus metrics using InstrumentHttpHandler.
func ExampleInstrumentHttpHandler() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	// Wrap with metrics middleware
	instrumented := prometrics.InstrumentHttpHandler("hello_handler", handler)

	// Simulate a request
	req := httptest.NewRequest("GET", "http://example.com/hello", nil)
	w := httptest.NewRecorder()
	instrumented.ServeHTTP(w, req)

	fmt.Println("Response Code:", w.Code)
	// Output:
	// Response Code: 200
}

// HealthMiddleware demonstrates how to register runtime metrics
// (goroutines, GC stats, memory usage, etc.) using HealthMiddleware.
func ExampleHealthMiddleware() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// Wrap with Go runtime collector
	instrumented := prometrics.HealthMiddleware(handler)

	req := httptest.NewRequest("GET", "http://example.com/health", nil)
	w := httptest.NewRecorder()
	instrumented.ServeHTTP(w, req)

	fmt.Println("Health endpoint responded:", w.Body.String())
	// Output:
	// Health endpoint responded: OK
}

// ExampleTrackCRUD desmonstrates how to use TrackCRUD function to track crud operation
// total and crud operation duration
func ExampleTrackCRUD() {
	// Simulate a "create" operation for "person" object.
	done := prometrics.TrackCRUD("person", "create")
	defer done(time.Now())

	// Perform your operation here...
	time.Sleep(120 * time.Millisecond)

	// Output:
	// (no output, metrics are exported to Prometheus)
}

// ExampleSetObjectCount demonstrates how to Set and increment/decrement object count
func ExampleSetObjectCount() {
	// Set initial count of "person" objects.
	prometrics.SetObjectCount("person", 42)

	// Increment and decrement.
	prometrics.IncObjectCount("person")
	prometrics.DecObjectCount("person")

	// Output:
	// (no output, metrics are exported to Prometheus)
}

// ExampleEndToEnd demonstrates using prometrics in a small end-to-end HTTP server.
// This example won’t run indefinitely in tests, but illustrates a real-world setup.
func Example_endToEnd() {
	mux := http.NewServeMux()
	mux.Handle("/person", prometrics.InstrumentHttpHandler("person_handler", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Person CRUD example")
	})))

	// Expose Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	srv := httptest.NewServer(mux)
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/person")
	fmt.Println("Status:", resp.StatusCode)

	// Simulate scraping metrics
	metricsResp, err := http.Get(srv.URL + "/metrics")
	if err != nil {
		return
	}

	defer metricsResp.Body.Close()
	fmt.Println("Metrics exposed:", metricsResp.StatusCode == 200)

	time.Sleep(100 * time.Millisecond)
	// Output:
	// Status: 200
	// Metrics exposed: true
}
