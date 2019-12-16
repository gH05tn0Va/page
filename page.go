package page

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

var client *http.Client

type DocCallBack func(*goquery.Document) interface{}

type Pages struct {
	Urls []string
}

type Job struct {
	Set       *Pages
	CallBack  DocCallBack
	OutFormat int
	Work      func(string) interface{}
	Out       struct {
		Slice          []string
		IntSlice       []int
		CustomSlice    []interface{}
		Map            map[string]string
		IntMap         map[string]int
		CustomMap      map[string]interface{}
		MapSlice       map[string][]string
		IntMapSlice    map[string][]int
		CustomMapSlice map[string][]interface{}
	}
}

func On(urls []string) *Pages {
	s := new(Pages)
	s.Urls = urls
	return s
}

func OnRange(format string, begin, end int) *Pages {
	s := new(Pages)
	for i := begin; i <= end; i++ {
		s.Urls = append(s.Urls, fmt.Sprintf(format, i))
	}
	return s
}

func (j *Job) setCallback() {
	j.Work = func(url string) interface{} {
		doc, err := GetPageBody(url)
		if err != nil {
			return err
		}
		if j.CallBack != nil {
			return j.CallBack(doc)
		}
		return doc
	}
}

func (s *Pages) Do(f DocCallBack) *Job {
	j := new(Job)
	j.OutFormat = StringOut
	j.Set = s
	j.CallBack = f
	j.setCallback()
	return j
}

func (j *Job) SetOutput(out interface{}) *Job {
	switch out.(type) {
	case []int:
		j.OutFormat = IntOut
		j.Out.IntSlice = out.([]int)
	case []string:
		j.OutFormat = StringOut
		j.Out.Slice = out.([]string)
	case map[string]int:
		j.OutFormat = IntMapOut
		j.Out.IntMap = out.(map[string]int)
	case map[string]string:
		j.OutFormat = StringMapOut
		j.Out.Map = out.(map[string]string)
	case []interface{}:
		j.OutFormat = InterfaceOut
		j.Out.CustomSlice = out.([]interface{})
	case map[string]interface{}:
		j.OutFormat = InterfaceMapOut
		j.Out.CustomMap = out.(map[string]interface{})
	case nil:
		break
	default:
		return nil
	}
	return j
}

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = &http.Client{Transport: tr}
}

func GetPageBody(url string) (*goquery.Document, error) {
	resp, err := client.Get(url)
	if err != nil {
		log.Print(err.Error())
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
