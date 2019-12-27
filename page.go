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

type (
	Doc            = *goquery.Document
	Urls           []string
	PageOutput     = []PageTaskOutput
	PageTaskOutput = [][]string
	PageTask       func(Doc) [][]string
)

/*
Output:
{
	"http://url1": []PageOutput{
		[]PageTaskOutput{
			[]string{
				"value1-1", // sub-task 1
				"value1-2", // sub-task 2
				...
			}, // Selector child 1
			[]string{
				"value2-1", // sub-task 1
				"value2-2", // sub-task 2
				...
			}, // Selector child 2
			...
		}, // task 1
		... // task 2
	},
	"http://url2": ...
	...
{
 */

type Job struct {
	Set        []string
	Tasks      []PageTask
	Output     map[string]PageOutput
	WorkerFunc func(string, *Job) error
}

func On(urls []string) *Urls {
	s := new(Urls)
	*s = append(*s, urls...)
	return s
}

func OnRange(format string, begin, end int) *Urls {
	s := new(Urls)
	for i := begin; i <= end; i++ {
		*s = append(*s, fmt.Sprintf(format, i))
	}
	return s
}

func (s *Urls) AddTask(f PageTask) *Job {
	j := new(Job)
	j.Set = *s
	j.Output = make(map[string]PageOutput)
	j.WorkerFunc = func(url string, j *Job) error {
		doc, err := GetPageBody(url)
		if err != nil {
			return err
		}
		var out []PageTaskOutput
		for _, taskFunc := range j.Tasks {
			v := taskFunc(doc)
			if v != nil {
				out = append(out, v)
			}
		}
		j.Output[url] = out
		return nil
	}
	return j.AddTask(f)
}

func (j *Job) AddTask(f PageTask) *Job {
	j.Tasks = append(j.Tasks, f)
	return j
}

func (s *Urls) Text() *Job {
	j := new(Job)
	j.Set = *s
	j.Output = make(map[string]PageOutput)
	j.WorkerFunc = func(url string, j *Job) error {
		doc, err := GetPageBody(url)
		if err != nil {
			return err
		}
		var out []PageTaskOutput
		for _, taskFunc := range j.Tasks {
			v := taskFunc(doc)
			if v != nil {
				out = append(out, v)
			}
		}
		j.Output[url] = out
		return nil
	}
	return j.Text()
}

func (j *Job) Text() *Job {
	j.Tasks = append(j.Tasks,
		func(doc Doc) [][]string{
			return [][]string{{doc.Text()}}
		})
	return j
}

func (w *Worker) ListAll() (out [][]string) {
	for _, v := range w.Wait().Job.Output {
		// PageTaskOutput
		for _, vv := range v {
			out = append(out, vv...)
		}
	}
	return
}

func (w *Worker) List(i int) (out []string) {
	for _, v := range w.Wait().Job.Output {
		// PageTaskOutput
		for _, vv := range v {
			for _, vvv := range vv {
				if i >= len(vvv) {
					continue
				}
				out = append(out, vvv[i])
			}
		}
	}
	return
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
