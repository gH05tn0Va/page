package page

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"sync"
)

type (
	Urls struct {
		Data []string
		Tag  string
	}
	PagingJob struct {
		BaseJob
		Output  OutMap
		TaskMap TaskMap
		Tasks   []PagingTask
	}

	Doc        = *goquery.Document
	PagingTask func(Doc) []MultiOut
)

func OnOne(url string) Urls {
	return OnMany([]string{url})
}

func OnMany(urls []string) Urls {
	return Urls{urls, ""}
}

func OnRange(format string, begin, end, step int) Urls {
	var s Urls
	for i := begin; i <= end; i += step {
		s.Data = append(s.Data, fmt.Sprintf(format, i))
	}
	return s
}

func New() *PagingJob {
	pj := new(PagingJob)
	pj.Output = Out    // Global output
	pj.TaskMap = Tasks // Global task map
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

		doc, err := getPageBody(url)
		if err != nil {
			log.Printf("GET %s ERR %s", url, err.Error())
			return err
		}
		if DebugWorker {
			log.Printf("GET %s OK", url)
		}

		for id, taskFunc := range pj.Tasks {
			v := taskFunc(doc)
			if DebugWorker {
				log.Printf("Task %d %s", id, url)
			}
			if v != nil {
				lock.Lock()
				pj.TaskMap[url] = append(pj.TaskMap[url],
					len(pj.Output[url]))
				pj.Output[url] = append(pj.Output[url], v...)
				lock.Unlock()
			}
		}
		lock.Lock()
		pj.TaskMap[url] = append(pj.TaskMap[url],
			len(pj.Output[url]))
		lock.Unlock()

		return nil
	}
}

func (s Urls) AddToWorker(w *Worker) *Worker {
	Tags.AddJob(s.Tag, w.Job.(*PagingJob))
	w.Job.(*PagingJob).Add(s)
	return w.Add(s.Data)
}

func (pj *PagingJob) Add(s Urls) {
	pj.Input = append(pj.Input, s.Data...)
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
		func(doc Doc) []MultiOut {
			return []MultiOut{{doc.Text()}}
		})
	return pj
}

func (pj *PagingJob) Out() OutMap {
	if pj.Worker != nil {
		return pj.Worker.Out()
	}
	return pj.Run().Out()
}

func (w *Worker) Out() OutMap {
	return w.Wait().Job.(*PagingJob).Output
}
