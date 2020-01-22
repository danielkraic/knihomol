package bookfinder

import (
	"context"
	"time"

	"github.com/danielkraic/knihomol/books"
	"github.com/danielkraic/knihomol/storage"
	log "github.com/sirupsen/logrus"
)

//BookFinder find books and its items in library
type BookFinder interface {
	//FindBooks finds books in library by given find query
	FindBooks(findQuery string) ([]*books.Book, error)

	//GetBook return book details
	GetBook(bookID string) (*books.Book, error)

	//FindBooksItems finds book items in library for given book
	FindBooksItems(bookID string) *books.Book

	// GetItemURL return URl to view book items
	GetItemURL(bookID string) string
}

//FindBooksItems finds books items in parallel
func FindBooksItems(finder BookFinder, booksToFind []*books.Book, webStorage *storage.Storage, storageTimeout time.Duration) {
	resultsChan := make(chan *books.Book, len(booksToFind))

	for _, bookToFind := range booksToFind {
		go func(bookToFind *books.Book) {
			log.Debugf("finding items for book %s STARTED", bookToFind.ID)

			result := finder.FindBooksItems(bookToFind.ID)
			result.LastUpdate = time.Now().Format("2006-01-02 15:04:05")
			log.Debugf("finding items for book %s DONE. items=%d. error=%s", bookToFind.ID, len(result.Items), result.Error)

			resultsChan <- result
		}(bookToFind)
	}

	for i := 0; i < len(booksToFind); i++ {
		result := <-resultsChan
		ctx, cancel := context.WithTimeout(context.Background(), storageTimeout)
		defer cancel()

		err := webStorage.SaveBook(ctx, result)
		if err != nil {
			log.Warnf("failed to save book %s: %s", result.ID, err)
		}
	}
}
