package handlers

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/danielkraic/knihomol/bookfinder"
	"github.com/danielkraic/knihomol/bookfinder/kjftt"
	"github.com/danielkraic/knihomol/storage"
)

type getBooksHandler struct {
	webStorage *storage.Storage
	finder     bookfinder.BookFinder
	timeout    time.Duration
}

//Book contains book details
type getBooksResult struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"url"`
}

// NewGetBooksHandler creates handler to get books
func NewGetBooksHandler(webStorage *storage.Storage, timeout time.Duration) http.Handler {
	return &getBooksHandler{
		webStorage: webStorage,
		finder:     kjftt.NewKJFTT(timeout),
		timeout:    timeout,
	}
}

// ServeHTTP serves http request
func (h *getBooksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	storedBooks, err := h.webStorage.GetBooks(ctx)
	if err != nil {
		log.Errorf("failed to get books: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := make([]getBooksResult, 0)
	for _, book := range storedBooks {
		result = append(result, getBooksResult{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
			URL:    h.finder.GetItemURL(book.ID),
		})
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Errorf("failed to encode json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
