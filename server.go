package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danielkraic/knihomol/configuration"
	"github.com/danielkraic/knihomol/controllers"
	"github.com/danielkraic/knihomol/middlewares"
	"github.com/danielkraic/knihomol/resources"
	"github.com/danielkraic/knihomol/views"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

	api := views.NewAPI(server.configuration, server.controller)

	basicAuth := middlewares.NewAuthenticationMiddleware(server.configuration.Auth.Username, server.configuration.Auth.Password)

	r.Use(forceSsl)
	r.HandleFunc("/", booksView.Index).Methods(http.MethodGet)
	r.HandleFunc("/list", booksView.ListBooks).Methods(http.MethodGet)

	restricted := r.PathPrefix("/restricted/").Subrouter()
	restricted.Use(basicAuth.Middleware)
	restricted.HandleFunc("/add-book", booksView.AddBook).Methods(http.MethodGet, http.MethodPost)
	restricted.HandleFunc("/remove-book", booksView.RemoveBook).Methods(http.MethodPost)
	restricted.HandleFunc("/refresh", booksView.Refresh).Methods(http.MethodPost)

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(basicAuth.Middleware)
	apiRouter.HandleFunc("/books", api.GetBooks).Methods(http.MethodGet)
	apiRouter.HandleFunc("/books/add", api.AddBook).Methods(http.MethodPost)
	apiRouter.HandleFunc("/books/remove", api.RemoveBook).Methods(http.MethodPost)
	apiRouter.HandleFunc("/books/remove-many", api.RemoveBooks).Methods(http.MethodPost)
	apiRouter.HandleFunc("/books/refresh", api.RefreshBooks).Methods(http.MethodPost)

	return r, nil
}

func forceSsl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("GO_ENV") == "production" {
			if r.Header.Get("x-forwarded-proto") != "https" {
				sslURL := "https://" + r.Host + r.RequestURI
				http.Redirect(w, r, sslURL, http.StatusTemporaryRedirect)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
