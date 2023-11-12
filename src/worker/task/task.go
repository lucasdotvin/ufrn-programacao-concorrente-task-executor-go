package task

import "time"

type Task struct {
	ID    int
	Cost  time.Duration
	Type  Type
	Value int
}

func NewTask(id int, cost time.Duration, taskType Type, value int) *Task {
	return &Task{
		ID:    id,
		Cost:  cost,
		Type:  taskType,
		Value: value,
	}
}
