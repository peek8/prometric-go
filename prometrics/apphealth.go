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
	AppUptime   = CreateGauge("app_uptime_seconds", "App uptime in seconds", nil)
	MemoryAlloc = CreateGauge("app_memory_alloc_bytes", "Memory allocated in bytes", nil)
	Goroutines  = CreateGauge("app_goroutines", "Current goroutines", nil)
	GCCount     = CreateCounter("app_gc_total", "Total garbage collections", nil)
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
// It should be called in a go routine, eg. if we want to collect metrics in 10 seconds interval, it should be called as follows:
// ctx, cancel := context.WithCancel(context.Background(), 10)
// go collectSystemMetricsLoop(ctx)
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

// HealthMiddleware for net/http
func HealthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		updateHealthMetrics()
		next.ServeHTTP(w, r)
	})
}
