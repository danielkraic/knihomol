package kjftt

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/danielkraic/knihomol/books"
)

//FindBooks finds books in library by given find query
func (kjftt *KJFTT) FindBooks(findQuery string) ([]*books.BookDetails, error) {
	doc, err := kjftt.httpGet(getSearchURL(findQuery))
	if err != nil {
		return nil, err
	}

	var result []*books.BookDetails

	doc.Find("li.record").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".title")
		author := s.Find(".author").Text()

		href, found := title.Attr("href")
		if !found {
			log.Printf("href attr not found")
			return
		}

		query := strings.Split(href, "?")
		if len(query) != 2 {
			log.Printf("failed to split href %s", href)
			return
		}

		items, err := url.ParseQuery(query[1])
		if err != nil {
			log.Printf("failed to parse href query. href=%s. query=%v. items=%v. error=%s", href, query, items, err)
			return
		}
		id := items.Get("id")

		result = append(result, &books.BookDetails{
			ID:     id,
			Title:  title.Text(),
			Author: author,
		})
	})

	return result, nil
}

func getSearchURL(findQuery string) string {
	values := url.Values{}
	values.Add("theme", "ttkjf")
	values.Add("term_1", findQuery)
	return "https://chamo.kis3g.sk/search/query?" + values.Encode()
}
