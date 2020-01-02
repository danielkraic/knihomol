package books

import (
	"fmt"
)

//BookItem contains book details and status
type BookItem struct {
	ItemID    string `json:"item_id" bson:"item_id"`
	Available bool   `json:"available" bson:"available"`
	Status    string `json:"status" bson:"status"`
	Location  string `json:"location" bson:"location"`
}

func (item BookItem) String() string {
	return fmt.Sprintf("itemid=%s. available=%v. status=%s. location=%s", item.ItemID, item.Available, item.Status, item.Location)
}
