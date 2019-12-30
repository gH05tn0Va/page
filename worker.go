package page

import (
	"log"
	"sync"
	"time"
)

var DebugWorker = false

type Worker struct {
	Job       Job
	Current   int
	Lock      sync.Mutex
}

type Job interface {
	Run() *Worker
	GetWorkFunc() func(string, Job, *sync.Mutex) error
	GetSet() []string
}

type BaseJob struct {
	Set        []string
	WorkerFunc func(string, Job, *sync.Mutex) error
}

func (j *BaseJob) GetWorkFunc() func(string, Job, *sync.Mutex) error {
	return j.WorkerFunc
}

func (j *BaseJob) GetSet() []string {
	return j.Set
}

func (j *BaseJob) Run() *Worker {
	w := new(Worker)
	w.Job = j
	return w.Run()
}

func (w *Worker) Run() *Worker {
	for _, s := range w.Job.GetSet() {
		w.Lock.Lock()
		w.Current++
		w.Debug("Added", s)
		w.Lock.Unlock()

		go func(s string) {
			err := w.Job.GetWorkFunc()(s, w.Job, &w.Lock)
			if err != nil {
				w.Warn("Error", err.Error())
			}
			w.Lock.Lock()
			w.Current--
			w.Debug("Finished", s)
			w.Lock.Unlock()
		}(s)
	}
	return w
}

func (w *Worker) Add(set []string) *Worker {
	for _, s := range set {
		w.Lock.Lock()
		w.Current++
		w.Debug("Added", s)
		w.Lock.Unlock()

		go func(s string) {
			err := w.Job.GetWorkFunc()(s, w.Job, &w.Lock)
			if err != nil {
				w.Warn("Error", err.Error())
			}
			w.Lock.Lock()
			w.Current--
			w.Debug("Finished", s)
			w.Lock.Unlock()
		}(s)
	}
	return w
}

func (w *Worker) AddOne(s string) *Worker {
	return w.Add([]string{s})
}

func (w *Worker) Wait() *Worker {
	for w.Current > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	return w
}

func (w *Worker) Warn(msg string, v interface{}) {
	log.Printf("[Pending: %d] %s: %v", w.Current, msg, v)
}

func (w *Worker) Debug(msg string, v interface{}) {
	if DebugWorker {
		w.Warn(msg, v)
	}
}
