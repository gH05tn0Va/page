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

func getResp(url string) (*http.Response, error) {
	resp, err := client.Get(url)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, errors.New(
			fmt.Sprintf("Status code: %d", resp.StatusCode))
	}
	return resp, nil
}

func getPageBody(url string) (Doc, error) {
	retry := 10
	retryWait := time.Second
	resp, err := getResp(url)

	for (err != nil) && retry > 0 {
		time.Sleep(retryWait)
		resp, err = getResp(url)

		retryWait *= 5
		retry--
	}

	if err != nil {
		log.Printf("FAIL %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
