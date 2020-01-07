package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/danielkraic/knihomol/bookfinder"
	"github.com/danielkraic/knihomol/bookfinder/kjftt"
	"github.com/danielkraic/knihomol/storage"
)

type findItemsHandler struct {
	webStorage *storage.Storage
	finder     bookfinder.BookFinder
	timeout    time.Duration
}

// NewFindItemsHandler creates handler to find books
func NewFindItemsHandler(webStorage *storage.Storage, timeout time.Duration) http.Handler {
	return &findItemsHandler{
		webStorage: webStorage,
		finder:     kjftt.NewKJFTT(timeout),
		timeout:    timeout,
	}
}

// ServeHTTP serves http request
func (h *findItemsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	booksToFind, err := h.webStorage.GetBooks(ctx)
	if err != nil {
		log.Errorf("failed to get books json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debugf("finding items for %d books", len(booksToFind))
	result := bookfinder.FindBooksItems(h.finder, booksToFind)

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Errorf("failed to encode json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
