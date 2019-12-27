package page

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	Job      *Job
	Current  int
	Lock     sync.Mutex
	Complete bool
}

func (j *Job) Run() *Worker {
	w := new(Worker)
	w.Job = j
	w.Current = len(j.Set)
	w.Complete = false
	for _, url := range j.Set {
		go func(url string) {
			err := j.WorkerFunc(url, j)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			}
			w.Lock.Lock()
			w.Current--
			w.Lock.Unlock()
		}(url)
	}
	return w
}

func (w *Worker) Wait() *Worker {
	for w.Current > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	w.Complete = true
	return w
}