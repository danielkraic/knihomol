package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danielkraic/knihomol/api/handlers"
	"github.com/danielkraic/knihomol/configuration"
	"github.com/danielkraic/knihomol/storage"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//API stores API configuration and storage collection
type API struct {
	configuration *configuration.Configuration
	storage       *storage.Storage
	version       *handlers.Version
}

//NewAPI creates API instance
func NewAPI(version *handlers.Version, apiConfiguration *configuration.Configuration, apiStorage *storage.Storage) (*API, error) {
	return &API{
		configuration: apiConfiguration,
		storage:       apiStorage,
		version:       version,
	}, nil
}

// Run starts API by running http server
func (apiInstance *API) Run(done chan os.Signal) {
	httpServer := &http.Server{
		Addr:    apiInstance.configuration.Addr,
		Handler: apiInstance.createRouter(),
	}

	go func() {
		log.Printf("Starting server on %s", apiInstance.configuration.Addr)
		log.Fatal(httpServer.ListenAndServe())
	}()

	<-done
	log.Println("Shutdown signal received. Exiting.")

	err := httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("failed to shutdown http server: %s", err)
	}
}

// CreateRouter creates http router
func (apiInstance *API) createRouter() *mux.Router {
	r := mux.NewRouter()

	versioned := func(route string) string {
		return fmt.Sprintf("%s%s", apiInstance.configuration.APIPrefix, route)
	}

	r.Handle(versioned("/books"), handlers.NewGetBooksHandler(apiInstance.storage, time.Duration(apiInstance.configuration.Timeout)*time.Second)).Methods("GET")
	r.Handle(versioned("/items"), handlers.NewFindItemsHandler(apiInstance.storage, time.Duration(apiInstance.configuration.Timeout)*time.Second)).Methods("GET")

	r.Handle("/version", handlers.NewVersionHandler(apiInstance.version)).Methods("GET")
	r.HandleFunc("/health", handlers.HealthHandlerFunc).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	return r
}
