// Base types for goroutine concurrent worker.
package page

import (
	"log"
	"sync"
	"time"
)

var DebugWorker = false

type WorkFunc func(string, Job, *sync.Mutex) error

type Worker struct {
	Job
	Current int
	Lock    sync.Mutex
}

type Job interface {
	Run() *Worker
	WorkFunc() WorkFunc
	GetInput() []string
}

type BaseJob struct {
	Worker *Worker
	Input  []string
}

func (j *BaseJob) Run() *Worker {
	j.Worker = new(Worker)
	j.Worker.Job = j
	return j.Worker.Run()
}

func (j *BaseJob) WorkFunc() WorkFunc {
	return nil
}

func (j *BaseJob) GetInput() []string {
	return j.Input
}

func Work(w *Worker, input []string) *Worker {
	for _, s := range input {
		w.Lock.Lock()
		w.Current++
		w.Debug("Added", s)
		w.Lock.Unlock()

		go func(s string) {
			err := w.Job.WorkFunc()(s, w.Job, &w.Lock)
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

func (w *Worker) Run() *Worker {
	return Work(w, w.Job.GetInput())
}

func (w *Worker) Add(input []string) *Worker {
	return Work(w, input)
}

func (w *Worker) AddOne(s string) *Worker {
	return w.Add([]string{s})
}

func (w *Worker) Wait() *Worker {
	for w.Current > 0 {
		time.Sleep(50 * time.Millisecond)
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
