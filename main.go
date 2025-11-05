package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"peek8.io/prometric-go/prometrics"
)

func main() {
	apiUsingGin()
}

func apiUsingGin() {
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
