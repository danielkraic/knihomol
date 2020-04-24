package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/danielkraic/knihomol/models"
	"github.com/danielkraic/knihomol/resources"
	"github.com/danielkraic/knihomol/resources/bookfinder"
	"github.com/danielkraic/knihomol/resources/bookfinder/kjftt"
	log "github.com/sirupsen/logrus"
)

//BooksController controller to manage books
type BooksController struct {
	storage *resources.Storage
	finder  bookfinder.BookFinder
	timeout time.Duration
}

//NewBooksController creates new BooksController
func NewBooksController(storage *resources.Storage, timeout time.Duration) *BooksController {
	return &BooksController{
		storage: storage,
		finder:  kjftt.NewKJFTT(timeout),
		timeout: timeout,
	}
}

// GetBookIDFromURL return bookID from book items URL
func (controller *BooksController) GetBookIDFromURL(url string) string {
	return controller.finder.GetBookIDFromURL(url)
}

//GetBooks returns all saved books
func (controller *BooksController) GetBooks(ctx context.Context) ([]*models.Book, error) {
	result, err := controller.storage.GetBooks(ctx)
	return result, err
}

//AddBook get book's details and available items add save book to storage
func (controller *BooksController) AddBook(ctx context.Context, bookID string) error {
	book, err := controller.finder.GetBook(bookID)
	if err != nil {
		return fmt.Errorf("add book '%s': %v", bookID, err)
	}
	return controller.storage.SaveBook(ctx, book)
}

//RemoveBook removes book from storage
func (controller *BooksController) RemoveBook(ctx context.Context, bookID string) error {
	return controller.storage.RemoveBook(ctx, bookID)
}

//Refresh refresh all books details and available items
func (controller *BooksController) Refresh(ctx context.Context) []error {
	books, err := controller.GetBooks(ctx)
	if err != nil {
		return []error{err}
	}

	resultsChan := make(chan *models.Book, len(books))

	for _, book := range books {
		go func(bookToFind *models.Book) {
			log.Debugf("refreshing book %s STARTED", bookToFind.ID)

			result, err := controller.finder.GetBook(bookToFind.ID)
			if err != nil {
				bookToFind.Error = err.Error()
				bookToFind.URL = controller.finder.GetItemURL(bookToFind.ID)
				bookToFind.Items = []*models.Item{}
				bookToFind.LastUpdate = time.Now().Format("2006-01-02 15:04:05")
				log.Warnf("finding items for book %s FAILED. error=%s", bookToFind.ID, bookToFind.Error)
				resultsChan <- bookToFind
			} else {
				result.LastUpdate = time.Now().Format("2006-01-02 15:04:05")
				log.Debugf("finding items for book %s DONE. items=%d", bookToFind.ID, len(result.Items))
				resultsChan <- result
			}
		}(book)
	}

	var result []error
	for i := 0; i < len(books); i++ {
		book := <-resultsChan
		err = controller.storage.SaveBook(ctx, book)
		if err != nil {
			result = append(result, err)
		}
	}
	return result
}
