package task

import "time"

type Task struct {
	Cost  time.Duration
	ID    uint32
	Value uint8
	Type  Type
}

func NewTask(id uint32, cost time.Duration, taskType Type, value uint8) *Task {
	return &Task{
		ID:    id,
		Cost:  cost,
		Type:  taskType,
		Value: value,
	}
}
