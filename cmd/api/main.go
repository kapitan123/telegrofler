package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/config"
	"github.com/kapitan123/telegrofler/internal/routes"

	log "github.com/sirupsen/logrus"
)

// Main entry point. Starts HTTP service
func main() {
	log.Info("Telegrofler: starting...")

	router := mux.NewRouter()

	app, err := routes.NewApp()

	if err != nil {
		log.Fatalf("The application could not be started: %v", err)
	}

	// AK TODO should close the whole api
	defer app.Close()

	app.AddRoutes(router)
	app.AddHandlers()

	log.Info("Telegrofler: listening on: ", config.ServerPort)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", config.ServerPort),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
