package ui

import (
	"fmt"
	"net/http"
	"text/template"

	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/danielkraic/knihomol/bookfinder"
	"github.com/danielkraic/knihomol/bookfinder/kjftt"
	"github.com/danielkraic/knihomol/storage"
)

type findItemsResultItem struct {
	BookID    string
	URL       string
	Author    string
	Title     string
	Available bool
	ItemID    string
	Location  string
	Status    string
	Error     string
}
type findItemsResult struct {
	Items []findItemsResultItem
}

type findItemsHandler struct {
	webStorage *storage.Storage
	finder     bookfinder.BookFinder
	timeout    time.Duration
	tmpl       *template.Template
}

// NewFindItemsHandler creates handler to find books
func NewFindItemsHandler(webStorage *storage.Storage, timeout time.Duration) (http.Handler, error) {
	tmpl, err := template.ParseGlob("./templates/*")
	if err != nil {
		return nil, fmt.Errorf("failed to parse temlates: %s", err)
	}

	return &findItemsHandler{
		webStorage: webStorage,
		finder:     kjftt.NewKJFTT(timeout),
		timeout:    timeout,
		tmpl:       tmpl,
	}, nil
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

	items := make([]findItemsResultItem, 0)
	for _, book := range result {
		if book.Error != nil {
			items = append(items, findItemsResultItem{
				BookID: book.Book.ID,
				URL:    book.Book.URL,
				Title:  book.Book.Title,
				Author: book.Book.Author,
				Error:  book.Error.Error(),
			})
			continue
		}

		for _, item := range book.Book.Items {
			items = append(items, findItemsResultItem{
				BookID:    book.Book.ID,
				URL:       book.Book.URL,
				Title:     book.Book.Title,
				Author:    book.Book.Author,
				Available: item.Available,
				ItemID:    item.ItemID,
				Location:  item.Location,
				Status:    item.Status,
			})
		}
	}

	h.tmpl.ExecuteTemplate(w, "items.html", &findItemsResult{
		Items: items,
	})
}
