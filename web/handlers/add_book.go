package handlers

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/danielkraic/knihomol/bookfinder"
	"github.com/danielkraic/knihomol/bookfinder/kjftt"
	"github.com/danielkraic/knihomol/books"
	"github.com/danielkraic/knihomol/storage"
)

type addBookRequest struct {
	ID string `json:"id"`
}

type addBookHandler struct {
	webStorage *storage.Storage
	finder     bookfinder.BookFinder
	timeout    time.Duration
}

// NewAddBookHandler creates handler to save book
func NewAddBookHandler(webStorage *storage.Storage, timeout time.Duration) http.Handler {
	return &addBookHandler{
		webStorage: webStorage,
		finder:     kjftt.NewKJFTT(timeout),
		timeout:    timeout,
	}
}

// ServeHTTP serves http request
func (h *addBookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	var bookID addBookRequest
	err := json.NewDecoder(r.Body).Decode(&bookID)
	if err != nil {
		log.Errorf("failed to decode request body: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := h.finder.GetBook(bookID.ID)
	if err != nil {
		log.Errorf("failed to get details for book with ID %s: %s", bookID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	toFind := []*books.Book{book}
	bookfinder.FindBooksItems(h.finder, toFind, h.webStorage, h.timeout)
}
