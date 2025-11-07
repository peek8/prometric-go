package prometrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		handler := c.FullPath()

		HttpRequestsInFlight.WithLabelValues(handler).Inc()
		defer HttpRequestsInFlight.WithLabelValues(handler).Dec()

		reqLength := c.Request.ContentLength

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()

		if handler == "" {
			handler = "unknown"
		}

		if reqLength > 0 {
			HttpRequestSize.WithLabelValues(handler, c.Request.Method, status).Observe(float64(reqLength))
		}

		HttpRequestsTotal.WithLabelValues(handler, c.Request.Method, status).Inc()
		HttpRequestDuration.WithLabelValues(handler, c.Request.Method, status).Observe(duration)
		HttpResponseSize.WithLabelValues(handler, c.Request.Method, status).Observe(float64(c.Writer.Size()))
	}
}

func GinHealthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		updateHealthMetrics()
		c.Next()
	}
}
