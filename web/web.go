package web

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"

	"github.com/danielkraic/knihomol/configuration"
	"github.com/danielkraic/knihomol/storage"
	"github.com/danielkraic/knihomol/web/handlers"
	"github.com/danielkraic/knihomol/web/ui"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//Web stores Web configuration and storage collection
type Web struct {
	configuration *configuration.Configuration
	storage       *storage.Storage
	version       *handlers.Version
}

//NewWeb creates Web instance
func NewWeb(version *handlers.Version, webConfiguration *configuration.Configuration, webStorage *storage.Storage) (*Web, error) {
	return &Web{
		configuration: webConfiguration,
		storage:       webStorage,
		version:       version,
	}, nil
}

// Run starts Web by running http server
func (webInstance *Web) Run(done chan os.Signal) {
	router, err := webInstance.createRouter()
	if err != nil {
		log.Fatalf("failed to web create router: %s", err)
	}

	httpServer := &http.Server{
		Addr:    webInstance.configuration.Addr,
		Handler: router,
	}

	go func() {
		log.Infof("Starting server on %s", webInstance.configuration.Addr)
		log.Fatal(httpServer.ListenAndServe())
	}()

	<-done
	log.Info("Shutdown signal received. Exiting.")

	err = httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("failed to shutdown http server: %s", err)
	}
}

// CreateRouter creates http router
func (webInstance *Web) createRouter() (*mux.Router, error) {
	r := mux.NewRouter()

	// UI
	getItemsHandler, err := ui.NewGetItemsHandler(webInstance.storage, time.Duration(webInstance.configuration.Timeout)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to create find items handler: %s", err)
	}
	r.Handle("/", getItemsHandler).Methods("GET")

	// API
	versioned := func(route string) string {
		return fmt.Sprintf("%s%s", webInstance.configuration.APIPrefix, route)
	}

	r.Handle(versioned("/books"), handlers.NewGetBooksHandler(webInstance.storage, time.Duration(webInstance.configuration.Timeout)*time.Second)).Methods("GET")
	r.Handle(versioned("/items"), handlers.NewFindItemsHandler(webInstance.storage, time.Duration(webInstance.configuration.Timeout)*time.Second)).Methods("GET")
	r.Handle(versioned("/save"), handlers.NewAddBookHandler(webInstance.storage, time.Duration(webInstance.configuration.Timeout)*time.Second)).Methods("POST")

	r.Handle("/version", handlers.NewVersionHandler(webInstance.version)).Methods("GET")
	r.HandleFunc("/health", handlers.HealthHandlerFunc).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	return r, nil
}
