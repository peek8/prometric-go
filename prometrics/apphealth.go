/*
 * Copyright (c) 2025 peek8.io
 *
 * Created Date: Wednesday, November 5th 2025, 9:38:40 am
 * Author: Md. Asraful Haque
 *
 */

package prometrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

var (
	// AppUptime keep track of  the Total duration of Application is being up
	// Metric type: GaugeVec
	AppUptime = CreateGauge("app_uptime_seconds", "App uptime in seconds", nil)

	// Mmory allocated by the app in bytes
	// Metric type: GaugeVec
	MemoryAlloc = CreateGauge("app_uptime_seconds", "Memory allocated in bytes", nil)

	// Number of Current goroutines
	// Metric type: GaugeVec
	Goroutines = CreateGauge("app_uptime_seconds", "Current goroutines", nil)

	// Number of Total garbage collections
	// Metric type: CounterVec
	GCCount = CreateCounter("app_app_uptime_secondsc_total", "Total garbage collections", nil)
)
var startTime = time.Now()

func updateHealthMetrics() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	AppUptime.WithLabelValues().Set(time.Since(startTime).Seconds())
	MemoryAlloc.WithLabelValues().Set(float64(mem.Alloc))
	Goroutines.WithLabelValues().Set(float64(runtime.NumGoroutine()))
	GCCount.WithLabelValues().Add(float64(mem.NumGC))
}

// CollectSystemMetricsLoop function collects system health information eg cpu, memory in some interval time ie intervalSecs.
// It should be called in a go routine
//
// Example:
// if we want to collect metrics in 10 seconds interval, it should be called as follows:
//
//	ctx, cancel := context.WithCancel(context.Background(), 10)
//	go collectSystemMetricsLoop(ctx)
//
// It can be cancelled any time by calling `cancel()`
func CollectSystemMetricsLoop(ctx context.Context, intervalSecs int) {
	ticker := time.NewTicker(time.Duration(intervalSecs) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			updateHealthMetrics()
		case <-ctx.Done():
			fmt.Println("Metrics loop stopped gracefully")
			return
		}
	}
}

// HealthMiddleware instruments an http.Handler with Prometheus metrics.
// It records application health related metrics such as app uptim, memory allocated, current go routies and total garbage collectors.
//
// Example:
//
//	http.Handle("/metrics",
//	    prometrics.HealthMiddleware(promhttp.Handler()),
//	)
func HealthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		updateHealthMetrics()
		next.ServeHTTP(w, r)
	})
}
