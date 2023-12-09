package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"task-executor/worker"
	"task-executor/worker/repository"
	"task-executor/worker/result"
	"task-executor/worker/service"
	"task-executor/worker/task"
	"time"
)

const dataFilePath = "./data/shared.txt"

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

func randomTask(id uint32, taskType task.Type) *task.Task {
	cost := time.Duration(rand.Float64() * 0.01 * float64(time.Millisecond))
	value := rand.Intn(11)

	return task.NewTask(id, cost, taskType, uint8(value))
}

func generateTasks(tasksCount int, t int, e int) []*task.Task {
	writingTasksCount := int(math.Ceil(float64(tasksCount) * (float64(e) / 100.0)))

	tasks := make([]*task.Task, tasksCount)

	for i := 0; i < writingTasksCount; i++ {
		tasks[i] = randomTask(uint32(i), task.Write)
	}

	for i := writingTasksCount; i < tasksCount; i++ {
		tasks[i] = randomTask(uint32(i), task.Read)
	}

	rand.Shuffle(tasksCount, func(i, j int) {
		tasks[i], tasks[j] = tasks[j], tasks[i]
		tasks[i].ID = uint32(i)
		tasks[j].ID = uint32(j)
	})

	return tasks
}

func main() {
	n, t, e, err := parseArgs()

	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	fmt.Println("Parâmetros recebidos:")
	fmt.Printf("n: %d, t: %d, e: %d\n", n, t, e)
	fmt.Printf("Arquivo de dados: %s\n", dataFilePath)

	tasksCount := int(math.Pow10(n))
	tasks := generateTasks(tasksCount, t, e)

	fmt.Printf("Total de tarefas: %d\n", tasksCount)
	fmt.Println("")

	tasksChan := make(chan *task.Task, tasksCount)
	resultsChan := make(chan *result.Result, tasksCount)

	repository, err := repository.NewRepository(dataFilePath)

	if err != nil {
		panic(err)
	}

	service := service.NewService(repository)

	processingStart := time.Now()

	fmt.Println("Iniciando workers...")

	workers := make([]*worker.Worker, t)

	for i := 0; i < t; i++ {
		workers[i] = worker.NewWorker(i, tasksChan, resultsChan, service)
		go workers[i].Start()
	}

	fmt.Println("Workers iniciados")
	fmt.Println("")
	fmt.Println("Iniciando distribuição de tarefas...")

	for _, t := range tasks {
		tasksChan <- t
	}

	close(tasksChan)

	fmt.Println("Tarefas distribuídas")
	fmt.Println("")
	fmt.Println("Aguardando resultados...")

	var resultsSample [10]*result.Result
	sampleIndexStep := tasksCount / 10
	sampleIndex := 0

	for i := 0; i < tasksCount; i++ {
		r := <-resultsChan

		if i == sampleIndex*sampleIndexStep {
			resultsSample[sampleIndex] = r
			sampleIndex++
		}
	}

	close(resultsChan)

	for _, r := range resultsSample {
		fmt.Printf("Resultado da tarefa %d: %d em %s\n", r.TaskID, r.Result, r.ElapsedTime)
	}

	fmt.Println("Resultados recebidos")
	fmt.Println("")

	processingDuration := time.Since(processingStart)

	fmt.Printf("Tempo de processamento total: %s\n", processingDuration)
}
