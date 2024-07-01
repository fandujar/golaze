package golaze

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
)

type WorkerConfig struct {
	Tasks           []*Task
	State           *State
	ConcurrentTasks int
}

type Worker struct {
	*WorkerConfig
	taskQueue *TaskQueue
	lock      sync.Mutex
	shutdown  chan bool
}

// NewWorker creates a new worker
func NewWorker(config *WorkerConfig) *Worker {
	if config.Tasks == nil {
		config.Tasks = make([]*Task, 0)
	}

	// default to 2 concurrent tasks
	if config.ConcurrentTasks == 0 {
		config.ConcurrentTasks = 2
	}

	if config.State == nil {
		config.State = &State{
			Data: make(map[string]interface{}),
		}
	}

	queue := &TaskQueue{
		tasks:   make([]Task, 0),
		enqueue: make(chan Task),
		dequeue: make(chan Task),
	}

	w := &Worker{
		WorkerConfig: config,
		taskQueue:    queue,
		shutdown:     make(chan bool),
		lock:         sync.Mutex{},
	}

	return w
}

// AddTask adds a task to the worker
func (w *Worker) AddTask(task *Task) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	// avoid adding the same task twice
	for _, t := range w.Tasks {
		if t.Name == task.Name {
			return fmt.Errorf("task %s already exists", task.Name)
		}
	}

	log.Debug().Msgf("adding task %s", task.Name)
	w.Tasks = append(w.Tasks, task)
	return nil
}

// RemoveTask removes a task from the worker
func (w *Worker) RemoveTask(task *Task) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	for i, t := range w.Tasks {
		if t.Name == task.Name {
			w.Tasks = append(w.Tasks[:i], w.Tasks[i+1:]...)
		}
	}

	return nil
}

// Start starts the worker
func (w *Worker) Start() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		s := <-signals
		log.Info().Msgf("received signal: %v", s)
		w.Stop()
	}()
}

// WaitShutdown waits for the worker to finish
func (w *Worker) WaitShutdown() {
	<-w.shutdown
}

// Stop stops the worker
func (w *Worker) Stop() {
	w.shutdown <- true
}
