// Helper functions for http requests
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

	client = &http.Client{
		Transport: tr,
		//Timeout:   time.Duration(10 * time.Second),
	}
}

func GetPageBody(url string) (Doc, error) {
	retry := 10
	retryWait := time.Second
	resp, err := client.Get(url)

	for (err != nil || resp.StatusCode != 200) && retry > 0 {
		time.Sleep(retryWait)
		if DebugWorker {
			log.Printf("RETRY [%d/3] %s", 4-retry, err.Error())
		}
		resp, err = client.Get(url)

		retryWait *= 5
		retry--
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
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
