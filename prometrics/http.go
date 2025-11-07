/*
 * Copyright (c) 2025 peek8.io
 *
 * Created Date: Monday, November 3rd 2025, 3:23:33 pm
 * Author: Md. Asraful Haque
 *
 */

package prometrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// type HttpMetricHandler struct {
// 	Name        HTTPMetricName
// 	HandlerFunc func(string, http.Handler) http.Handler
// 	Metric      prometheus.Collector
// }

// var (
// 	HttpMetricHandlers []HttpMetricHandler
// )

// InstrumentHttpHandler instruments an http.Handler with Prometheus metrics.
//
// It records total requests, request duration, in-flight requests, request size and response size.
// The "handler" label is set from the given handlerName parameter.
// The returned handler can be used directly in an http.ServeMux.
//
// Example:
//
//	http.Handle("/api",
//	    prometrics.InstrumentHttpHandler("api", myHandler),
//	)
func InstrumentHttpHandler(handlerName string, next http.Handler) http.Handler {
	h := promhttp.InstrumentHandlerInFlight(HttpRequestsInFlight.WithLabelValues(handlerName),
		promhttp.InstrumentHandlerDuration(
			HttpRequestDuration.MustCurryWith(handlerLabel(handlerName)),
			promhttp.InstrumentHandlerCounter(
				HttpRequestsTotal.MustCurryWith(handlerLabel(handlerName)),
				promhttp.InstrumentHandlerResponseSize(
					HttpResponseSize.MustCurryWith(handlerLabel(handlerName)),
					promhttp.InstrumentHandlerRequestSize(
						HttpRequestSize.MustCurryWith(handlerLabel(handlerName)),
						next,
					),
				),
			),
		),
	)

	return h
}

// func init() {
// 	reqTotal := HttpMetricHandler{
// 		Name: HttpRequestsTotalMetric,
// 		HandlerFunc: func(handlerName string, next http.Handler) http.Handler {
// 			return promhttp.InstrumentHandlerCounter(
// 				HttpRequestsTotal.MustCurryWith(handlerLabel(handlerName)), next)
// 		},
// 		Metric: HttpRequestsTotal,
// 	}

// 	reqDuration := HttpMetricHandler{
// 		Name: HttpRequestDurationMetric,
// 		HandlerFunc: func(handlerName string, next http.Handler) http.Handler {
// 			return promhttp.InstrumentHandlerDuration(
// 				HttpRequestDuration.MustCurryWith(handlerLabel(handlerName)), next)
// 		},
// 		Metric: HttpRequestDuration,
// 	}

// 	requestInFlight := HttpMetricHandler{
// 		Name: HttpRequestsInFlightMetric,
// 		HandlerFunc: func(handlerName string, next http.Handler, p prometheus.Collector) http.Handler {
// 			return promhttp.InstrumentHandlerInFlight(p, next)
// 		},
// 		Metric: HttpRequestsInFlight,
// 	}

// 	handlers := []HttpMetricHandler{reqTotal, reqDuration, requestInFlight}

// 	// lo.Reduce(handlers, func(agg http.Handler, mh MetricHandler, _ int) http.Handler{
// 	// 	return mh.HandlerFunc(handlerName, agg)
// 	// }, next)
// }

func handlerLabel(name string) prometheus.Labels {
	return prometheus.Labels{"path": name}
}

// HttpMiddleware is a generic version to wrap muxes or routers easily
func HttpMiddleware(next http.Handler) http.Handler {
	return InstrumentHttpHandler("default", next)
}
