
[![Go Reference](https://pkg.go.dev/badge/github.com/peek8/prometric-go.svg)](https://pkg.go.dev/github.com/peek8/prometric-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/peek8/prometric-go?v2)](https://goreportcard.com/report/github.com/peek8/prometric-go)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache%20License%202.0-blue)](LICENSE)

# prometric-go
**_An open source project from [Peek8.io](https://peek8.io/)._** 

Prometric-go is a lightweight and configurable **Prometheus instrumentation library** for Go web applications.

`prometric-go` makes it easy to add standardized **HTTP**, **health**, and **business-level metrics** to your services — with minimal code and full Prometheus compatibility.



## ✨ Features

### Plug-and-play HTTP Middleware 
This middleware can be used with any `net/http` or Gin handler. It Automatically exposes:
   - `http_requests_total{path,method,code}`-  Total number of HTTP requests processed
   - `http_request_duration_seconds{path,method,code}` - Histogram of HTTP request durations in seconds.
   - `http_in_flight_requests{path}` - Number of HTTP requests currently being handled.
   - `http_request_size_bytes{path,method,code}` - Size of incoming HTTP requests in bytes.
   - `http_response_size_bytes{path,method,code}` - Size of outgoing HTTP responses in bytes.

### Plug-and-play Health Middleware
This middleware can be used with any `net/http` or Gin handler. It Automatically exposes:
- `app_uptime_seconds` - App uptime in seconds.    
- `app_cpu_usage_percent` - CPU usage (percent) of GO process.    
- `app_allocated_memory` - Memory allocated in bytes.    
- `app_go_routines` - Number of Current goroutines.   
- `app_garbage_collections_count` - Total completed garbage collections count.    

### CRUD Monitoring Functions
Exposes some utility functions to track crud operation and business object metrics. With these functions, the following metrics can be exposed:
- `crud_operations_total{"object", "operation"}` - Total CRUD operations.
- `object_operation_duration_seconds{"object", "operation"}` - CRUD duration.
- `object_count{"object"}` -  Current number of objects.

### 💡 100% compatible with Prometheus + Grafana


## 🧩 Installation

```bash
go get github.com/peek8/prometric-go/prometrics
```


## 🚀 Quick Start

### Using `net/http`
If you are using net/http to expose your api, you can use the Middleware to expose the necessary HTTP metrics and App Health Metrics. Use crud metric functions (eg. prometrics.TrackCRUD) to expose metrics for crud Operations.

```Go
package main

import (
	"net/http"
    "time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/peek8/prometric-go/prometrics"
)

func createPerson(w http.ResponseWriter, r *http.Request) {
    // Track the person create operation
    defer prometrics.TrackCRUD("person", "Create")(time.Now())
    prometrics.IncObjectCount("person")

    // add some delay to imitate the real creation operation
	time.Sleep(200 * time.Millisecond)

    w.Write([]byte("Person created"))
}

func main() {
	mux := http.NewServeMux()

	// Wrap your business handler
	mux.Handle("/person", prometrics.InstrumentHttpHandler("person_handler",
		http.HandlerFunc(createPerson)))

	// Use HealthMiddleware at /metrics endpoint
	mux.Handle("/metrics", prometrics.HealthMiddleware(promhttp.Handler()))

	http.ListenAndServe(":8080", mux)
}
```

### Using `Gin`
If you are using gin, you can use the the Gin Middlewares from prometric library ie `prometrics.GinMiddleware()` for http metrics and `prometrics.GinHealthMiddleware` for app health. Use the CRUD tracking functions as before.

```Go
package main

import (
    "time"

    "github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/peek8/prometric-go/prometrics"
)

func main() {
	r := gin.Default()
	r.Use(prometrics.GinMiddleware())
	r.Use(prometrics.GinHealthMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})

	r.GET("/person", func(c *gin.Context) {
		defer prometrics.TrackCRUD("person", "Get")(time.Now())

		c.JSON(200, gin.H{"name": "asraf"})
	})

	r.POST("/person", func(c *gin.Context) {
		defer prometrics.TrackCRUD("person", "create")(time.Now())
		prometrics.IncObjectCount("person")

		time.Sleep(200 * time.Millisecond)
		c.JSON(201, gin.H{"status": "created"})
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.Run(":7080")
}
```

The health middleware is used to collect app health info if some request are done in some endpoint (eg. at /metric handler). If you want to collect the health info at some interval irrespective of the some endpoint request, you can use `prometrics.CollectSystemMetricsLoop()` function. It should be called in a go routine. if we want to collect metrics in 10 seconds interval, it should be called as follows:

```Go
ctx, cancel := context.WithCancel(context.Background(), 10)
go prometrics.CollectSystemMetricsLoop(ctx)
// It can be cancelled any time by calling `cancel()`
```


## 📚 Documentation

Full API reference available at:
👉 [pkg.go.dev/github.com/peek8/prometric-go/prometrics](https://pkg.go.dev/github.com/peek8/prometric-go/prometrics)

## Prometric GO Sample
If you want to explore how this library can be used to collect metrics in Prometheus, and visualize everything in Grafana — including a load-test using k6 to generate traffic.
Please have a look at the sample github repo:

🔗 GitHub Repo: https://github.com/peek8/prometric-go-sample

## 📜 License
- Apache 2.0, see more details at [LICENSE File](./LICENSE).

## Community
Prometric-go is a [Peek8](https://peek8.io/) open source project.
Learn about our open source work and portfolio [here](https://peek8.io/#products).
If you want to collaborate with us or Invest at [Peek8](https://peek8.io/), please [contact us here](https://peek8.io/#contact).

Issues and PRs are most welcome! Whether it's docs, code improvement, or examples — contributions help the community.

Last but not the least, If you like this project and it seems to be helpful to you, please consider giving the repository a ⭐. Your support helps build a better developer-first observability ecosystem.
