package prometrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		HttpRequestsInFlight.Inc()
		defer HttpRequestsInFlight.Dec()

		if c.Request.ContentLength > 0 {
			HttpRequestSize.WithLabelValues(c.FullPath(), c.Request.Method).Observe(float64(c.Request.ContentLength))
		}

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()
		handler := c.FullPath()
		if handler == "" {
			handler = "unknown"
		}

		HttpRequestsTotal.WithLabelValues(handler, c.Request.Method, status).Inc()
		HttpRequestDuration.WithLabelValues(handler, c.Request.Method).Observe(duration)
		HttpResponseSize.WithLabelValues(handler, c.Request.Method).Observe(float64(c.Writer.Size()))
	}
}

func GinHealthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		updateHealthMetrics()
		c.Next()
	}
}
