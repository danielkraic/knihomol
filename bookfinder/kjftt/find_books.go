package kjftt

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/danielkraic/knihomol/books"
)

//FindBooks finds books in library by given find query
func (kjftt *KJFTT) FindBooks(findQuery string) ([]*books.Book, error) {
	doc, err := kjftt.httpGet(getSearchURL(findQuery))
	if err != nil {
		return nil, err
	}

	var result []*books.Book

	doc.Find("li.record").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".title")
		author := s.Find(".author").Text()

		href, found := title.Attr("href")
		if !found {
			log.Warn("href attr not found")
			return
		}

		query := strings.Split(href, "?")
		if len(query) != 2 {
			log.Warnf("failed to split href %s", href)
			return
		}

		items, err := url.ParseQuery(query[1])
		if err != nil {
			log.Warnf("failed to parse href query. href=%s. query=%v. items=%v. error=%s", href, query, items, err)
			return
		}
		id := items.Get("id")

		result = append(result, &books.Book{
			ID:     id,
			Title:  title.Text(),
			Author: author,
			URL:    kjftt.GetItemURL(id),
		})
	})

	return result, nil
}
