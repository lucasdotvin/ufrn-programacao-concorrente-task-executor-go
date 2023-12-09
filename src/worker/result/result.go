package result

import "time"

type Result struct {
	TaskID      uint32
	Result      int
	ElapsedTime time.Duration
}

func NewResult(taskID uint32, result int, elapsedTime time.Duration) *Result {
	return &Result{
		TaskID:      taskID,
		Result:      result,
		ElapsedTime: elapsedTime,
	}
}
