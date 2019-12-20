package kjftt

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/danielkraic/knihomol/books"
)

//FindBooksItem finds book items in library for given book
func (kjftt *KJFTT) FindBooksItem(book *books.BookDetails) ([]*books.BookItem, error) {
	doc, err := kjftt.httpGet(getItemURL(book))
	if err != nil {
		return nil, err
	}

	var result []*books.BookItem

	doc.Find("#tabContents-1 tbody").Each(func(ibody int, body *goquery.Selection) {
		body.Find("tr").Each(func(itr int, tr *goquery.Selection) {
			// if tr.("display: none")
			ulozenie := tr.Find("td:nth-child(1)").Text()
			signatura := tr.Find("td:nth-child(2)").Text()
			ciarovyKod := tr.Find("td:nth-child(3)").Text()
			dostupnost := tr.Find("td:nth-child(4)").Text()
			status := tr.Find("td:nth-child(5)").Text()

			result = append(result, &books.BookItem{
				Details:   book,
				ID:        ciarovyKod,
				Available: false,
				Location:  fmt.Sprintf("u=%s s=%s d=%s s=%s", ulozenie, signatura, dostupnost, status),
			})
		})
	})

	return result, nil
}

func getItemURL(book *books.BookDetails) string {
	values := url.Values{}
	values.Add("theme", "ttkjf")
	values.Add("id", book.ID)
	return "https://chamo.kis3g.sk/lib/item?" + values.Encode()
}
