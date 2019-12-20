package books

import (
	"fmt"
	"strings"
)

//BookDetails contains book details
type BookDetails struct {
	ID     string `json:"book_id" bson:"book_id"`
	Title  string `json:"title" bson:"title"`
	Author string `json:"author" bson:"author"`
}

func (book BookDetails) String() string {
	items := strings.Split(book.Title, "/")
	titleSimple := strings.TrimSpace(items[0])

	items = strings.Split(items[1], ";")
	authorSimple := strings.TrimSpace(items[0])

	return fmt.Sprintf("id=%s. title=%s. author=%s", book.ID, titleSimple, authorSimple)
}
