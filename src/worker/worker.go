package worker

import (
	"task-executor/worker/reading"
	"task-executor/worker/result"
	"task-executor/worker/task"
	"task-executor/worker/writing"
	"time"
)

type Worker struct {
	ID              int
	Tasks           <-chan *task.Task
	Results         chan<- *result.Result
	ReadRepository  *reading.Repository
	WriteRepository *writing.Repository
}

func NewWorker(
	id int,
	tasks <-chan *task.Task,
	result chan<- *result.Result,
	readRepository *reading.Repository,
	writeRepository *writing.Repository,
) *Worker {
	return &Worker{
		ID:              id,
		Tasks:           tasks,
		Results:         result,
		ReadRepository:  readRepository,
		WriteRepository: writeRepository,
	}
}

func (w *Worker) Start() {
	for t := range w.Tasks {
		switch t.Type {
		case task.Read:
			w.Results <- w.read(t)
		case task.Write:
			w.Results <- w.write(t)
		}
	}
}

func (w *Worker) read(t *task.Task) *result.Result {
	start := time.Now()
	time.Sleep(t.Cost)

	r, err := w.ReadRepository.Read()

	if err != nil {
		panic(err)
	}

	return &result.Result{
		TaskID:      t.ID,
		Result:      r,
		ElapsedTime: time.Since(start),
	}
}

func (w *Worker) write(t *task.Task) *result.Result {
	start := time.Now()
	time.Sleep(t.Cost)

	current, err := w.ReadRepository.Read()

	if err != nil {
		panic(err)
	}

	r := current + t.Value
	err = w.WriteRepository.Write(r)

	if err != nil {
		panic(err)
	}

	return result.NewResult(
		t.ID,
		r,
		time.Since(start),
	)
}
