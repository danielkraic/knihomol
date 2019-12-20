package bookfinder

import (
	"github.com/danielkraic/knihomol/books"
)

//BookFinder find books and its items in library
type BookFinder interface {
	//FindBooks finds books in library by given find query
	FindBooks(findQuery string) ([]*books.BookDetails, error)

	//FindBooksItem finds book items in library for given book
	FindBooksItem(book *books.BookDetails) ([]*books.BookItem, error)
}

//FindBooksItemResult result of one find o book items
type FindBooksItemResult struct {
	items []*books.BookItem
	err   error
}

//FindBooksItem finds books items in parallel
func FindBooksItem(finder BookFinder, booksToFind []*books.BookDetails) []*FindBooksItemResult {
	resultsChan := make(chan *FindBooksItemResult, len(booksToFind))

	for _, book := range booksToFind {
		go func(book *books.BookDetails) {
			items, err := finder.FindBooksItem(book)
			resultsChan <- &FindBooksItemResult{
				items: items,
				err:   err,
			}
		}(book)
	}

	var results []*FindBooksItemResult
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}
