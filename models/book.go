package models

//Book contains book details
import (
	"fmt"
	"strings"
)

//Book contains book details
type Book struct {
	ID          string  `json:"id" bson:"_id"`
	Title       string  `json:"title" bson:"title"`
	Author      string  `json:"author" bson:"author"`
	Description string  `json:"description" bson:"description"`
	Publisher   string  `json:"publisher" bson:"publisher"`
	URL         string  `json:"url" bson:"url"`
	Items       []*Item `json:"items" bson:"items"`
	LastUpdate  string  `json:"last_update" bson:"last_update"`
	Error       string  `json:"error" bson:"error"`
}

func (book Book) String() string {
	items := strings.Split(book.Title, "/")
	titleSimple := strings.TrimSpace(items[0])

	items = strings.Split(items[1], ";")
	authorSimple := strings.TrimSpace(items[0])

	return fmt.Sprintf("id=%s. title=%s. author=%s", book.ID, titleSimple, authorSimple)
}

//Item contains single book item in library
type Item struct {
	ItemID    string `json:"item_id" bson:"id"`
	Available bool   `json:"available" bson:"available"`
	Status    string `json:"status" bson:"status"`
	Location  string `json:"location" bson:"location"`
}

func (item Item) String() string {
	return fmt.Sprintf("id=%s. available=%v. status=%s. location=%s", item.ItemID, item.Available, item.Status, item.Location)
}
