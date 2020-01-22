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

type getItemsResultItem struct {
	BookID     string
	URL        string
	Author     string
	Title      string
	Available  bool
	ItemID     string
	Location   string
	Status     string
	Error      string
	LastUpdate string
}
type getItemsResult struct {
	Items []getItemsResultItem
}

type getItemsHandler struct {
	webStorage *storage.Storage
	finder     bookfinder.BookFinder
	timeout    time.Duration
	tmpl       *template.Template
}

// NewGetItemsHandler creates handler to find books
func NewGetItemsHandler(webStorage *storage.Storage, timeout time.Duration) (http.Handler, error) {
	tmpl, err := template.ParseGlob("./templates/*")
	if err != nil {
		return nil, fmt.Errorf("failed to parse temlates: %s", err)
	}

	return &getItemsHandler{
		webStorage: webStorage,
		finder:     kjftt.NewKJFTT(timeout),
		timeout:    timeout,
		tmpl:       tmpl,
	}, nil
}

// ServeHTTP serves http request
func (h *getItemsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	result, err := h.webStorage.GetBooks(ctx)
	if err != nil {
		log.Errorf("failed to get books json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items := make([]getItemsResultItem, 0)
	for _, book := range result {
		if book.Error != "" {
			items = append(items, getItemsResultItem{
				BookID:     book.ID,
				URL:        book.URL,
				Title:      book.Title,
				Author:     book.Author,
				LastUpdate: book.LastUpdate,
				Error:      book.Error,
			})
			continue
		}

		for _, item := range book.Items {
			items = append(items, getItemsResultItem{
				BookID:     book.ID,
				URL:        book.URL,
				Title:      book.Title,
				Author:     book.Author,
				Available:  item.Available,
				ItemID:     item.ItemID,
				Location:   item.Location,
				Status:     item.Status,
				LastUpdate: book.LastUpdate,
				Error:      book.Error,
			})
		}
	}

	err = h.tmpl.ExecuteTemplate(w, "items.html", &getItemsResult{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
