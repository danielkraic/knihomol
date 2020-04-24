package views

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/danielkraic/knihomol/configuration"
	"github.com/danielkraic/knihomol/controllers"
)

//BooksView renders books UI
type BooksView struct {
	controller *controllers.BooksController
	timeout    time.Duration
	tmpl       *template.Template
}

//NewBooksView creates BooksView
func NewBooksView(conf *configuration.Configuration, controller *controllers.BooksController) (*BooksView, error) {
	tmpl, err := template.ParseGlob("./templates/*")
	if err != nil {
		return nil, fmt.Errorf("failed to parse temlates: %w", err)
	}

	timeoutSec := time.Duration(conf.Timeout) * time.Second

	return &BooksView{
		controller: controller,
		timeout:    timeoutSec,
		tmpl:       tmpl,
	}, nil
}

//Index renders index view
func (booksView BooksView) Index(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), booksView.timeout)
	defer cancel()

	books, err := booksView.controller.GetBooks(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var items []indexResultItem
	for _, book := range books {
		if book.Error != "" {
			items = append(items, indexResultItem{
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
			items = append(items, indexResultItem{
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

	err = booksView.tmpl.ExecuteTemplate(w, "index.html", &indexResult{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type indexResult struct {
	Items []indexResultItem
}

type indexResultItem struct {
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

//AddBook renders view to add book
func (booksView *BooksView) AddBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), booksView.timeout)
	defer cancel()

	var result manageBookResult
	result.Success = false

	if r.Method == http.MethodPost {
		bookID, err := booksView.getBookID(r)
		if err != nil {
			result.Err = err
		} else {
			result.BookID = bookID
			err = booksView.controller.AddBook(ctx, bookID)
			if err != nil {
				result.Err = err
			} else {
				result.Success = true
			}
		}
	}

	err := booksView.tmpl.ExecuteTemplate(w, "add-book.html", result)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (booksView *BooksView) getBookID(r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", fmt.Errorf("parse form: %w", err)
	}
	bookID := r.PostForm.Get("bookid")
	if bookID == "" {
		return "", fmt.Errorf("empty Book ID in form")
	}

	if strings.HasPrefix(bookID, "http") {
		url := bookID
		bookID = booksView.controller.GetBookIDFromURL(url)
		if bookID == "" {
			return "", fmt.Errorf("invalid url '%s'", url)
		}
	}

	return bookID, nil
}

type manageBookResult struct {
	BookID  string
	Success bool
	Err     error
}

//RemoveBook removes book
func (booksView *BooksView) RemoveBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), booksView.timeout)
	defer cancel()

	var result manageBookResult
	result.Success = false

	bookID, err := booksView.getBookID(r)
	if err != nil {
		result.Err = err
	} else {
		result.BookID = bookID
		err = booksView.controller.RemoveBook(ctx, bookID)
		if err != nil {
			result.Err = err
		} else {
			result.Success = true
		}
	}

	err = booksView.tmpl.ExecuteTemplate(w, "remove-book.html", result)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
