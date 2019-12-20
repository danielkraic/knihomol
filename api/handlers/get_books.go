package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/danielkraic/knihomol/bookfinder"
	"github.com/danielkraic/knihomol/bookfinder/kjftt"
	"github.com/danielkraic/knihomol/storage"
)

type getBooksHandler struct {
	apiStorage *storage.Storage
	finder     bookfinder.BookFinder
	timeout    time.Duration
}

// NewGetBooksHandler creates Version handler
func NewGetBooksHandler(apiStorage *storage.Storage, timeout time.Duration) http.Handler {
	return &getBooksHandler{
		apiStorage: apiStorage,
		finder:     kjftt.NewKJFTT(timeout),
		timeout:    timeout,
	}
}

// ServeHTTP serves http request
func (h *getBooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	result, err := h.apiStorage.GetBooks(ctx)
	if err != nil {
		log.Printf("failed to get books: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("failed to encode json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
