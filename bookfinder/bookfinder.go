package bookfinder

import (
	"github.com/danielkraic/knihomol/books"
	log "github.com/sirupsen/logrus"
)

//BookFinder find books and its items in library
type BookFinder interface {
	//FindBooks finds books in library by given find query
	FindBooks(findQuery string) ([]*books.Book, error)

	//GetBook return book details
	GetBook(bookID string) (*books.Book, error)

	//FindBooksItems finds book items in library for given book
	FindBooksItems(bookID string) (*books.Book, error)

	// GetItemURL return URl to view book items
	GetItemURL(bookID string) string
}

//FindBooksItemsResult result of one find o book items
type FindBooksItemsResult struct {
	Book  *books.Book `json:"book"`
	Error error       `json:"error"`
}

//FindBooksItems finds books items in parallel
func FindBooksItems(finder BookFinder, booksToFind []*books.Book) []*FindBooksItemsResult {
	resultsChan := make(chan *FindBooksItemsResult, len(booksToFind))

	for _, bookToFind := range booksToFind {
		go func(bookToFind *books.Book) {
			log.Debugf("finding items for book %s STARTED", bookToFind.ID)

			result, err := finder.FindBooksItems(bookToFind.ID)
			log.Debugf("finding items for book %s DONE. items=%d. error=%s", bookToFind.ID, len(result.Items), err)

			resultsChan <- &FindBooksItemsResult{
				Book:  result,
				Error: err,
			}
		}(bookToFind)
	}

	var results []*FindBooksItemsResult
	for i := 0; i < len(booksToFind); i++ {
		result := <-resultsChan
		results = append(results, result)
	}

	return results
}
