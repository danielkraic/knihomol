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

type findItemsHandler struct {
	apiStorage *storage.Storage
	finder     bookfinder.BookFinder
	timeout    time.Duration
}

// NewFindItemsHandler creates Version handler
func NewFindItemsHandler(apiStorage *storage.Storage, timeout time.Duration) http.Handler {
	return &findItemsHandler{
		apiStorage: apiStorage,
		finder:     kjftt.NewKJFTT(timeout),
		timeout:    timeout,
	}
}

// ServeHTTP serves http request
func (h *findItemsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	booksToFind, err := h.apiStorage.GetBooks(ctx)
	if err != nil {
		log.Printf("failed to get books json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := bookfinder.FindBooksItem(h.finder, booksToFind)

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("failed to encode json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
