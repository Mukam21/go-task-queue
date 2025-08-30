package worker

import (
	"math"
	"math/rand"
	"project_go-task-queue/pkg/models"
	"sync"
	"time"
)

type TaskQueue struct {
	Queue   chan *models.Task
	States  sync.Map
	Workers int
	wg      sync.WaitGroup
}

func NewTaskQueue(workers, queueSize int) *TaskQueue {
	return &TaskQueue{
		Queue:   make(chan *models.Task, queueSize),
		Workers: workers,
	}
}

func (tq *TaskQueue) Start() {
	for i := 0; i < tq.Workers; i++ {
		tq.wg.Add(1)
		go tq.worker(i)
	}
}

func (tq *TaskQueue) Stop() {
	close(tq.Queue)
	tq.wg.Wait()
}

func (tq *TaskQueue) worker(id int) {
	defer tq.wg.Done()
	for task := range tq.Queue {
		tq.processTask(task)
	}
}

func (tq *TaskQueue) processTask(task *models.Task) {
	task.Status = "running"
	tq.States.Store(task.ID, task)

	time.Sleep(time.Duration(100+rand.Intn(400)) * time.Millisecond)

	if rand.Intn(100) < 20 { // 20% ошибка
		task.Attempts++
		if task.Attempts > task.MaxRetries {
			task.Status = "failed"
			tq.States.Store(task.ID, task)
			return
		}
		backoff := time.Millisecond * time.Duration(100*math.Pow(2, float64(task.Attempts)))
		jitter := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(backoff + jitter)
		tq.processTask(task)
		return
	}

	task.Status = "done"
	tq.States.Store(task.ID, task)
}
