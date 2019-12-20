package kjftt

import (
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func (kjftt *KJFTT) httpGet(getURL string) (*goquery.Document, error) {
	log.Printf("\nGET %s\n", getURL)

	res, err := kjftt.client.Get(getURL)
	if err != nil {
		return nil, fmt.Errorf("http request failed. url=%s. error=%s", getURL, err)
	}

	log.Printf("%d\n\n", res.StatusCode)

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
