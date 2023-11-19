package service

import (
	"sync"
	"task-executor/worker/repository"
	"task-executor/worker/result"
	"task-executor/worker/task"
	"time"
)

type Service struct {
	repository *repository.Repository
	mutex      *sync.RWMutex
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		repository: repository,
		mutex:      new(sync.RWMutex),
	}
}

func (s *Service) Process(t *task.Task) *result.Result {
	switch t.Type {
	case task.Read:
		return s.read(t)
	case task.Write:
		return s.write(t)
	}

	panic("invalid task type")
}

func (s *Service) read(t *task.Task) *result.Result {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	start := time.Now()
	time.Sleep(t.Cost)

	r, err := s.repository.Read()

	if err != nil {
		panic(err)
	}

	return &result.Result{
		TaskID:      t.ID,
		Result:      r,
		ElapsedTime: time.Since(start),
	}
}

func (s *Service) write(t *task.Task) *result.Result {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	start := time.Now()
	time.Sleep(t.Cost)

	current, err := s.repository.Read()

	if err != nil {
		panic(err)
	}

	r := current + t.Value
	err = s.repository.Write(r)

	if err != nil {
		panic(err)
	}

	return result.NewResult(
		t.ID,
		r,
		time.Since(start),
	)
}
