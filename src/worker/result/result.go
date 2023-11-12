package result

import "time"

type Result struct {
	TaskID      int
	Result      int
	ElapsedTime time.Duration
}

func NewResult(taskID int, result int, elapsedTime time.Duration) *Result {
	return &Result{
		TaskID:      taskID,
		Result:      result,
		ElapsedTime: elapsedTime,
	}
}
