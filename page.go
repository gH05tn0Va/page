package page

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"sync"
)

var WorkerMap = make(map[*PagingJob]*Worker)

type (
	Urls struct {
		data []string
		tag  string
	}

	PagingJob struct {
		BaseJob
		PagingOutput
		Status int
		Tasks  []PagingTask
		Tags   map[string]string
	}

	Doc               = *goquery.Document
	OutputWithTag     []string
	OutputListWithTag []OutputWithTag
	PagingOutput      map[string]OutputListWithTag
	PagingTask        func(Doc) []OutputWithTag
)

func OnOne(url string) *Urls {
	return OnMany([]string{url})
}

func OnMany(urls []string) *Urls {
	s := new(Urls)
	s.data = append(s.data, urls...)
	return s
}

func OnRange(format string, begin, end, step int) *Urls {
	s := new(Urls)
	for i := begin; i <= end; i += step {
		s.data = append(s.data, fmt.Sprintf(format, i))
	}
	return s
}

func (s *Urls) PageTag(pageTag string) *Urls {
	s.tag = pageTag
	return s
}

func New() *PagingJob {
	pj := new(PagingJob)
	pj.Status = 0 // Not run
	pj.Tags = make(map[string]string)
	pj.PagingOutput = make(PagingOutput)
	return pj
}

func (pj *PagingJob) Run() *Worker {
	pj.Status = 1 // Running
	w := new(Worker)
	w.Job = pj
	WorkerMap[pj] = w
	return w.Run()
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

		tag := pj.Tags[url]
		if tag == "" {
			tag = url
		}

		for i, taskFunc := range pj.Tasks {
			v := taskFunc(doc)
			if v != nil {
				lock.Lock()
				pj.PagingOutput[tag] = append(pj.PagingOutput[tag], v...)
				lock.Unlock()
			}
			if DebugWorker {
				log.Printf("TASK %d %s OK", i, url)
			}
		}

		return nil
	}
}

func (s *Urls) AddToWorker(w *Worker) *Worker {
	w.Job.(*PagingJob).Add(s)
	return w.Add(s.data)
}

func (pj *PagingJob) Add(s *Urls) {
	pj.Input = append(pj.Input, s.data...)
	if s.tag != "" {
		JobMap[s.tag] = pj
		for _, url := range s.data {
			pj.Tags[url] = s.tag
		}
	}
}

func (s *Urls) AddTask(f PagingTask) *PagingJob {
	pj := New()
	pj.Add(s)
	return pj.AddTask(f)
}

func (pj *PagingJob) AddTask(f PagingTask) *PagingJob {
	pj.Tasks = append(pj.Tasks, f)
	return pj
}

func (s *Urls) Text() Job {
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

/*
func (pj *PagingJob) Get(s string) OutputListWithTag {
	return pj.PagingOutput[s]
}

func (pj *PagingJob) GetFirst(s string) OutputWithTag {
	return pj.Get(s).First()
}

func (w *Worker) Get(s string) OutputListWithTag {
	return w.Out()[s]
}

func (w *Worker) GetFirst(s string) OutputWithTag {
	return w.Get(s).First()
}
*/

func (pj *PagingJob) Out() PagingOutput {
	w, ok := WorkerMap[pj]
	if ok {
		return w.Out()
	}
	return pj.Run().Out()
}

func (w *Worker) Out() PagingOutput {
	return w.Wait().Job.(*PagingJob).PagingOutput
}
