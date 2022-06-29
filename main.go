package main

import (
	"HTTPChiSqlite/api"
	"HTTPChiSqlite/db"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func health(w http.ResponseWriter, r *http.Request) {
	if db.Instance == nil || db.Instance.GetDB() == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Timeout(30 * time.Second))

	// Default routes
	r.Get("/", root)
	r.Get("/health", health)

	// People controller
	r.Route("/people", api.PeopleRouter)

	port := 3000
	envPort, hasPort := os.LookupEnv("HTTP_PORT")
	if hasPort {
		outPort, err := strconv.Atoi(envPort)
		if err != nil {
			log.Fatal(err)
		} else if outPort < 1025 || outPort > 65535 {
			log.Fatal("port should be between 1025 and 65535")
		} else {
			port = outPort
		}
	}

	// Database implementation
	db.Instance = &db.SqliteDatabase{}
	db.Instance.InitDB()
	defer func(database db.Database) {
		err := database.CloseDB()
		if err != nil {
			log.Fatal(err)
		}
	}(db.Instance)

	// Gracefully exit on keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(port), r)
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	log.Print("Ready to serve requests on port " + strconv.Itoa(port))
	<-c
	log.Print("Gracefully shutting down")
	os.Exit(0)
}
