package page

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"time"
)

var client *http.Client

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = &http.Client{Transport: tr}
}

func GetPageBody(url string) (Doc, error) {
	resp, err := client.Get(url)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	retry := 3
	retryWait := time.Second
	for resp.StatusCode != 200 && retry > 0 {
		if DebugWorker {
			log.Printf("RETRY [%d/3] GET %s %d",
				4-retry, url, resp.StatusCode)
		}

		resp, err = client.Get(url)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}

		time.Sleep(retryWait)
		retryWait *= 5
		retry--
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(
			fmt.Sprintf("Status code: %d", resp.StatusCode))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
