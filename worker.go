package page

import (
	"sync"
	"time"
)

const (
	IntOut          = 1 // [int int ...]
	StringOut       = 2 // [string string ...]
	IntMapOut       = 3 // {url: int, url: int, ...}
	StringMapOut    = 4 // {url: int, url: int, ...}
	SubIntOut       = 5 // [int int ...]
	SubStringOut    = 6 // [string string ...]
	SubIntMapOut    = 7 // [url: [int ...], ...]
	SubStringMapOut = 8 // [url: [string ...], ...]

	InterfaceOut       = 9  // [interface{} ...]
	InterfaceMapOut    = 10 // [url: interface{}, ...]
	SubInterfaceOut    = 11 // [interface{} ...]
	SubInterfaceMapOut = 12 // [url: [interface{} ...], ...]
)

type Worker struct {
	Job     *Job
	Current int
	Lock    sync.Mutex
}

func (j *Job) Run() *Worker {
	w := new(Worker)
	w.Job = j
	w.Current = len(j.Set.Urls)
	for _, url := range j.Set.Urls {
		go func(url string) {
			out := j.Work(url)
			if out != nil {
				w.output(url, out)
			}
			w.Lock.Lock()
			w.Current--
			w.Lock.Unlock()
		}(url)
	}
	return w
}

func (j *Job) RunWithRetry(retry int) *Worker {
	job := new(Worker)
	job.Current = len(j.Set.Urls)
	for _, url := range j.Set.Urls {
		go func(url string) {
			doc, err := GetPageBody(url)
			for retry > 0 && err != nil {
				doc, err = GetPageBody(url)
			}
			if err == nil && j.CallBack != nil {
				out := j.CallBack(doc)
				if out != nil {
					job.Lock.Lock()
					if out != nil {
						job.output(url, out)
					}
					job.Lock.Unlock()
				}
			}
			job.Lock.Lock()
			job.Current--
			job.Lock.Unlock()
		}(url)
	}
	return job
}

func (w *Worker) output(url string, out interface{}) {
	w.Lock.Lock()
	j := w.Job
	switch j.OutFormat {
	case IntOut:
		j.Out.IntSlice = append(j.Out.IntSlice, out.(int))
	case StringOut:
		j.Out.Slice = append(j.Out.Slice, out.(string))
	case IntMapOut:
		j.Out.IntMap[url] = out.(int)
	case StringMapOut:
		j.Out.Map[url] = out.(string)
	case SubIntOut:
		for _, v := range out.([]interface{}) {
			j.Out.IntSlice = append(j.Out.IntSlice, v.(int))
		}
	case SubStringOut:
		for _, v := range out.([]interface{}) {
			j.Out.Slice = append(j.Out.Slice, v.(string))
		}
	case SubIntMapOut:
		for _, v := range out.([]interface{}) {
			j.Out.IntMapSlice[url] = append(j.Out.IntMapSlice[url], v.(int))
		}
	case SubStringMapOut:
		for _, v := range out.([]interface{}) {
			j.Out.MapSlice[url] = append(j.Out.MapSlice[url], v.(string))
		}
	case InterfaceOut:
		j.Out.CustomSlice = append(j.Out.CustomSlice, out)
	case InterfaceMapOut:
		j.Out.CustomMap[url] = out
	case SubInterfaceOut:
		for _, v := range out.([]interface{}) {
			j.Out.CustomSlice = append(j.Out.CustomSlice, v)
		}
	case SubInterfaceMapOut:
		for _, v := range out.([]interface{}) {
			j.Out.CustomMapSlice[url] = append(j.Out.CustomMapSlice[url], v)
		}
	}
	w.Lock.Unlock()
}

func (w *Worker) Wait() *Worker {
	for w.Current > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	return w
}

func (w *Worker) Collect() []interface{} {
	return w.Job.Out.CustomSlice
}