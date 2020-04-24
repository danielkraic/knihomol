package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/danielkraic/knihomol/configuration"
	"github.com/danielkraic/knihomol/controllers"
	"github.com/danielkraic/knihomol/middlewares"
	"github.com/danielkraic/knihomol/resources"
	"github.com/danielkraic/knihomol/views"
	"github.com/gorilla/mux"
)

//Server serves html requests
type Server struct {
	configuration *configuration.Configuration
	controller    *controllers.BooksController
}

//NewServer creates Server
func NewServer(conf *configuration.Configuration, storage *resources.Storage) (*Server, error) {
	return &Server{
		configuration: conf,
		controller:    controllers.NewBooksController(storage, time.Duration(conf.Timeout)*time.Second),
	}, nil
}

// Run starts http server
func (server *Server) Run(done chan os.Signal) {
	router, err := server.createRouter()
	if err != nil {
		log.Fatalf("create router: %s", err)
	}

	httpServer := &http.Server{
		Addr:    server.configuration.Addr,
		Handler: router,
	}

	go func() {
		log.Infof("starting server on %s", server.configuration.Addr)
		log.Fatal(httpServer.ListenAndServe())
	}()

	<-done
	log.Info("shutdown signal received. Exiting.")

	err = httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("shutdown http server: %s", err)
	}
}

// createRouter creates http router
func (server *Server) createRouter() (*mux.Router, error) {
	r := mux.NewRouter()

	booksView, err := views.NewBooksView(server.configuration, server.controller)
	if err != nil {
		return nil, fmt.Errorf("create BooksView: %w", err)
	}

	basicAuth := middlewares.NewAuthenticationMiddleware(server.configuration.Auth.Username, server.configuration.Auth.Password)

	r.HandleFunc("/", booksView.Index).Methods(http.MethodGet)

	restricted := r.PathPrefix("/restricted/").Subrouter()
	restricted.Use(basicAuth.Middleware)
	restricted.HandleFunc("/add-book", booksView.AddBook).Methods(http.MethodGet, http.MethodPost)
	restricted.HandleFunc("/remove-book", booksView.RemoveBook).Methods(http.MethodPost)

	return r, nil
}
