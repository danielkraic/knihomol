package kjftt

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/danielkraic/knihomol/books"
	log "github.com/sirupsen/logrus"
	"strings"
)

//FindBooksItems finds book items in library for given book
func (kjftt *KJFTT) FindBooksItems(bookID string) *books.Book {
	doc, err := kjftt.httpGet(kjftt.GetItemURL(bookID))
	if err != nil {
		return &books.Book{
			ID:    bookID,
			Error: err.Error(),
		}
	}

	var items []*books.BookItem

	doc.Find("#tabContents-1 tbody").Each(func(ibody int, body *goquery.Selection) {
		body.Find("tr").Each(func(itr int, tr *goquery.Selection) {
			style, found := tr.Attr("style")
			if found && style == "display: none;" {
				// skip hidden <tr>
				return
			}

			ulozenie := strings.TrimSpace(tr.Find("td:nth-child(1)").Text())
			signatura := strings.TrimSpace(tr.Find("td:nth-child(2)").Text())
			ciarovyKod := strings.TrimSpace(tr.Find("td:nth-child(3)").Text())
			dostupnost := strings.TrimSpace(tr.Find("td:nth-child(4)").Text())
			status := strings.TrimSpace(tr.Find("td:nth-child(5)").Text())

			items = append(items, &books.BookItem{
				ItemID:    ciarovyKod,
				Available: strings.HasPrefix(dostupnost, "Vypožičateľné") && strings.HasPrefix(status, "Dostupné"),
				Status:    fmt.Sprintf("%s %s", dostupnost, status),
				Location:  fmt.Sprintf("%s %s", ulozenie, signatura),
			})
		})
	})

	log.Debugf("items found: %d", len(items))

	return &books.Book{
		ID:     bookID,
		Title:  doc.Find(".title").First().Text(),
		Author: doc.Find(".author").First().Text(),
		URL:    kjftt.GetItemURL(bookID),
		Items:  items,
	}
}
