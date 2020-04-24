package bookfinder

import (
	"github.com/danielkraic/knihomol/models"
)

//BookFinder find books and its items in library
type BookFinder interface {
	//FindBooks finds books in library by given find query
	FindBooks(findQuery string) ([]*models.Book, error)

	//GetBook return book details and available items
	GetBook(bookID string) (*models.Book, error)

	// GetItemURL return URL to view book items
	GetItemURL(bookID string) string

	// GetBookIDFromURL return bookID from book items URL
	GetBookIDFromURL(url string) string
}
