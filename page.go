package page

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"sync"
)

type (
	PagingJob struct {
		BaseJob
		Output  PagingOutput
		Tasks   []PagingTask
	}

	Urls       []string
	Doc        = *goquery.Document
	PagingTask func(Doc) []OutputWithTag
)

func OnOne(url string) Urls {
	return OnMany([]string{url})
}

func OnMany(urls []string) Urls {
	return urls
}

func OnRange(format string, begin, end, step int) Urls {
	var s Urls
	for i := begin; i <= end; i += step {
		s = append(s, fmt.Sprintf(format, i))
	}
	return s
}

func New() *PagingJob {
	pj := new(PagingJob)
	pj.Output = Out // Global output
	return pj
}

func (pj *PagingJob) Run() *Worker {
	pj.Worker = new(Worker)
	pj.Worker.Job = pj
	return pj.Worker.Run()
}

func (pj *PagingJob) WorkFunc() WorkFunc {
	return func(url string, j Job, lock *sync.Mutex) error {
		pj, ok := j.(*PagingJob)
		if !ok {
			log.Fatalf("Expected *page.PagingJob but got %T", j)
		}

		doc, err := GetPageBody(url)
		if err != nil {
			log.Printf("GET %s ERR %s", url, err.Error())
			return err
		}
		if DebugWorker {
			log.Printf("GET %s OK", url)
		}

		for _, taskFunc := range pj.Tasks {
			v := taskFunc(doc)
			if v != nil {
				lock.Lock()
				pj.Output[url] = v
				lock.Unlock()
			}
		}

		return nil
	}
}

func (s Urls) AddToWorker(w *Worker) *Worker {
	w.Job.(*PagingJob).Add(s)
	return w.Add(s)
}

func (pj *PagingJob) Add(s Urls) {
	pj.Input = append(pj.Input, s...)
}

func (s Urls) AddTask(f PagingTask) *PagingJob {
	pj := New()
	pj.Add(s)
	return pj.AddTask(f)
}

func (pj *PagingJob) AddTask(f PagingTask) *PagingJob {
	pj.Tasks = append(pj.Tasks, f)
	return pj
}

func (s Urls) Text() Job {
	pj := New()
	pj.Add(s)
	return pj.Text()
}

func (pj *PagingJob) Text() *PagingJob {
	pj.Tasks = append(pj.Tasks,
		func(doc Doc) []OutputWithTag {
			return []OutputWithTag{{doc.Text()}}
		})
	return pj
}

func (pj *PagingJob) Out() PagingOutput {
	if pj.Worker != nil {
		return pj.Worker.Out()
	}
	return pj.Run().Out()
}

func (w *Worker) Out() PagingOutput {
	return w.Wait().Job.(*PagingJob).Output
}