package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"strconv"
	"sync"
	"task-executor/worker"
	"task-executor/worker/reading"
	"task-executor/worker/result"
	"task-executor/worker/task"
	"task-executor/worker/writing"
	"time"
)

const (
	dataDir    = "./data"
	resultsDir = "./resultados"
)

func parseArgs() (n int, t int, e int, err error) {
	if len(os.Args) != 4 {
		return 0, 0, 0, errors.New("expected exactly 3 arguments")
	}

	rawN := os.Args[1]
	n, err = strconv.Atoi(rawN)

	if err != nil {
		return 0, 0, 0, errors.New("expected a number for the first argument")
	}

	rawT := os.Args[2]
	t, err = strconv.Atoi(rawT)

	if err != nil {
		return 0, 0, 0, errors.New("expected a number for the second argument")
	}

	rawE := os.Args[3]
	e, err = strconv.Atoi(rawE)

	if err != nil {
		return 0, 0, 0, errors.New("expected a number for the third argument")
	}

	if e < 0 || e > 100 {
		return 0, 0, 0, errors.New("expected the third argument to be between 0 and 100")
	}

	return n, t, e, nil
}

func randomBool(truePercentage int) bool {
	return rand.Intn(100) < truePercentage
}

func randomTask(id int, writingPercentage int) *task.Task {
	isWrite := randomBool(writingPercentage)

	var taskType task.Type

	if isWrite {
		taskType = task.Write
	} else {
		taskType = task.Read
	}

	cost := time.Duration(rand.Float64() * 0.01 * float64(time.Second))
	value := rand.Intn(11)

	return task.NewTask(id, cost, taskType, value)
}

func main() {
	n, t, e, err := parseArgs()

	if err != nil {
		panic(err)
	}

	executionTimestamp := time.Now().Format("2006-01-02T15-04-05")
	dataFilePath := path.Join(dataDir, executionTimestamp)

	resultsFilePath := path.Join(resultsDir, executionTimestamp)
	resultsFile, err := os.OpenFile(resultsFilePath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(resultsFile, "n: %d, t: %d, e: %d\n", n, t, e)

	tasksCount := int(math.Pow10(n))
	tasks := make([]*task.Task, tasksCount)

	for i := 0; i < tasksCount; i++ {
		tasks[i] = randomTask(i, e)
	}

	fmt.Fprintf(resultsFile, "tasks: %d\n", tasksCount)

	tasksChan := make(chan *task.Task, tasksCount)
	resultsChan := make(chan *result.Result, tasksCount)

	fileMux := new(sync.RWMutex)

	writingRepository := writing.NewRepository(dataFilePath, fileMux)
	readingRepository := reading.NewRepository(dataFilePath, fileMux)

	workers := make([]*worker.Worker, t)

	for i := 0; i < t; i++ {
		workers[i] = worker.NewWorker(i, tasksChan, resultsChan, readingRepository, writingRepository)
		go workers[i].Start()
	}

	processingStart := time.Now()

	for _, t := range tasks {
		tasksChan <- t
	}

	close(tasksChan)

	for i := 0; i < tasksCount; i++ {
		r := <-resultsChan
		fmt.Fprintf(resultsFile, "task: %d, value: %d, duration: %s\n", r.TaskID, r.Result, r.ElapsedTime)
	}

	close(resultsChan)

	processingDuration := time.Since(processingStart)

	fmt.Fprintf(resultsFile, "processing duration: %s\n", processingDuration)

	if err := resultsFile.Close(); err != nil {
		panic(err)
	}
}
