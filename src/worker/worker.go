package worker

import (
	"task-executor/worker/result"
	"task-executor/worker/service"
	"task-executor/worker/task"
)

type Worker struct {
	id      int
	tasks   <-chan *task.Task
	results chan<- *result.Result
	service *service.Service
}

func NewWorker(
	id int,
	tasks <-chan *task.Task,
	results chan<- *result.Result,
	service *service.Service,
) *Worker {
	return &Worker{
		id:      id,
		tasks:   tasks,
		results: results,
		service: service,
	}
}

func (w *Worker) Start() {
	for t := range w.tasks {
		w.results <- w.service.Process(t)
	}
}
