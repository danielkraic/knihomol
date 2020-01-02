package books

import (
	"fmt"
	"strings"
)

//Book contains book details
type Book struct {
	ID     string      `json:"id" bson:"_id"`
	Title  string      `json:"title" bson:"title"`
	Author string      `json:"author" bson:"author"`
	URL    string      `json:"url" bson:"-"`
	Items  []*BookItem `json:"items" bson:"-"`
}

func (book Book) String() string {
	items := strings.Split(book.Title, "/")
	titleSimple := strings.TrimSpace(items[0])

	items = strings.Split(items[1], ";")
	authorSimple := strings.TrimSpace(items[0])

	return fmt.Sprintf("id=%s. title=%s. author=%s", book.ID, titleSimple, authorSimple)
}
