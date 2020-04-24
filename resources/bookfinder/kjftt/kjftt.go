package kjftt

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/danielkraic/knihomol/models"
	log "github.com/sirupsen/logrus"
)

//KJFTT finder for KJFTT
type KJFTT struct {
	client *http.Client
}

//NewKJFTT creates new KJFTT finder
func NewKJFTT(timeout time.Duration) *KJFTT {
	return &KJFTT{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// GetItemURL return URL to view book items
func (kjftt *KJFTT) GetItemURL(bookID string) string {
	values := url.Values{}
	values.Add("theme", "ttkjf")
	values.Add("id", bookID)
	return "https://chamo.kis3g.sk/lib/item?" + values.Encode()
}

// GetBookIDFromURL return bookID from book items URL
func (kjftt *KJFTT) GetBookIDFromURL(url string) string {
	r := regexp.MustCompile(`\Wid=[^&]+`)
	result := r.FindString(url)
	if result != "" {
		result = result[4:]
	}
	return result
}

func (kjftt *KJFTT) httpGet(getURL string) (*goquery.Document, error) {
	log.Debugf("\nGET %s", getURL)

	res, err := kjftt.client.Get(getURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed. url=%s. error=%s", getURL, err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("http request failed. url=%s. status=%d", getURL, res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read html response from url %s: %s", getURL, err)
	}

	return doc, nil
}

func getSearchURL(findQuery string) string {
	values := url.Values{}
	values.Add("theme", "ttkjf")
	values.Add("term_1", findQuery)
	return "https://chamo.kis3g.sk/search/query?" + values.Encode()
}

//FindBooks finds books in library by given find query
func (kjftt *KJFTT) FindBooks(findQuery string) ([]*models.Book, error) {
	doc, err := kjftt.httpGet(getSearchURL(findQuery))
	if err != nil {
		return nil, err
	}

	var result []*models.Book

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

		result = append(result, &models.Book{
			ID:     id,
			Title:  title.Text(),
			Author: author,
			URL:    kjftt.GetItemURL(id),
		})
	})

	return result, nil
}

//GetBook finds book details and available items
func (kjftt *KJFTT) GetBook(bookID string) (*models.Book, error) {
	url := kjftt.GetItemURL(bookID)
	doc, err := kjftt.httpGet(url)
	if err != nil {
		return nil, err
	}

	var items []*models.Item

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

			items = append(items, &models.Item{
				ItemID:    ciarovyKod,
				Available: strings.HasPrefix(dostupnost, "Vypožičateľné") && strings.HasPrefix(status, "Dostupné"),
				Status:    fmt.Sprintf("%s %s", dostupnost, status),
				Location:  fmt.Sprintf("%s %s", signatura, ulozenie),
			})
		})
	})

	log.Debugf("items found: %d", len(items))

	return &models.Book{
		ID:     bookID,
		Title:  parseTitle(doc.Find(".title").First().Text()),
		Author: parseAuthor(doc.Find(".author").First().Text()),
		URL:    url,
		Items:  items,
	}, nil
}

func parseTitle(text string) string {
	return strings.Split(text, "/")[0]
}

func parseAuthor(text string) string {
	items := strings.Split(text, ",")
	if len(items) > 2 {
		return fmt.Sprintf("%s %s", items[1], items[0])
	}
	return text
}
