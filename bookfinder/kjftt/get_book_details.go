package kjftt

import (
	"github.com/danielkraic/knihomol/books"
)

//GetBook return book details
func (kjftt *KJFTT) GetBook(bookID string) (*books.Book, error) {

	doc, err := kjftt.httpGet(kjftt.GetItemURL(bookID))
	if err != nil {
		return nil, err
	}

	title := doc.Find(".title").First().Text()
	author := doc.Find(".author").First().Text()

	return &books.Book{
		ID:     bookID,
		Title:  title,
		Author: author,
		URL:    kjftt.GetItemURL(bookID),
	}, nil
}
