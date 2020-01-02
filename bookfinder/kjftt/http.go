package kjftt

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func (kjftt *KJFTT) httpGet(getURL string) (*goquery.Document, error) {
	log.Debugf("\nGET %s", getURL)

	res, err := kjftt.client.Get(getURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed. url=%s. error=%s", getURL, err)
	}

	log.Debugf("%d", res.StatusCode)

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
		return nil, fmt.Errorf("failed to read html response from url %s: %s", getURL, err)
	}

	return doc, nil
}

func getSearchURL(findQuery string) string {
	values := url.Values{}
	values.Add("theme", "ttkjf")
	values.Add("term_1", findQuery)
	return "https://chamo.kis3g.sk/search/query?" + values.Encode()
}

// GetItemURL return URl to view book items
func (kjftt *KJFTT) GetItemURL(bookID string) string {
	values := url.Values{}
	values.Add("theme", "ttkjf")
	values.Add("id", bookID)
	return "https://chamo.kis3g.sk/lib/item?" + values.Encode()
}
