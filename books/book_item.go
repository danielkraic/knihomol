package books

import (
	"fmt"
)

//BookItem contains book details and status
type BookItem struct {
	Details   *BookDetails `json:"details" bson:"details"`
	ID        string       `json:"item_id" bson:"item_id"`
	Available bool         `json:"available" bson:"available"`
	Status    string       `json:"status" bson:"status"`
	Location  string       `json:"location" bson:"location"`
}

func (item BookItem) String() string {
	return fmt.Sprintf("%s. itemid=%s. available=%v. status=%s. location=%s", item.Details, item.ID, item.Available, item.Status, item.Location)
}
