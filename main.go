package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"peek8.io/prometric-go/prometrics"
)

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	persons   = make(map[int]Person)
	personMux sync.Mutex
	nextID    = 1
)

func main() {
	//apiUsingGin()
	httpApi()
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	personMux.Lock()
	defer personMux.Unlock()
	defer prometrics.TrackCRUD("person", "Create")(time.Now())

	var p Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.ID = nextID
	nextID++
	persons[p.ID] = p
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func getPersons(w http.ResponseWriter, r *http.Request) {
	personMux.Lock()
	defer personMux.Unlock()
	defer prometrics.TrackCRUD("persons", "Get")(time.Now())

	list := make([]Person, 0, len(persons))
	for _, p := range persons {
		list = append(list, p)
	}
	json.NewEncoder(w).Encode(list)
}

func httpApi() {
	// On test add one person
	persons[nextID] = Person{
		ID:   nextID,
		Name: "asraf",
	}
	nextID++

	r := mux.NewRouter()

	r.Handle("/person", prometrics.InstrumentHttpHandler("/person", http.HandlerFunc(createPerson))).Methods("POST")
	r.Handle("/persons", prometrics.InstrumentHttpHandler("/persons", http.HandlerFunc(getPersons))).Methods("GET")

	// expose metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	fmt.Println("Server listening on :7080")
	serverAddr := "0.0.0.0:7080"
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
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
