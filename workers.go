package golaze

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

type WorkerConfig struct {
	Tasks           []*Task
	ConcurrentTasks int
}

type Worker struct {
	*WorkerConfig
}

type WorkerServer struct {
	taskQueue *TaskQueue
	state     *State
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

	w := &Worker{
		WorkerConfig: config,
	}

	return w
}

// NewWorkerServer creates a new worker server
func NewWorkerServer() *WorkerServer {
	taskQueue := &TaskQueue{
		enqueue: make(chan Task, 10),
		dequeue: make(chan Task, 10),
	}

	state := &State{
		Data: make(map[string]interface{}),
	}

	return &WorkerServer{
		taskQueue: taskQueue,
		state:     state,
		lock:      sync.Mutex{},
		shutdown:  make(chan bool, 1),
	}
}

// AddTask adds a task to the worker
func (w *WorkerServer) AddTask(task *Task) error {
	w.taskQueue.enqueue <- *task

	return nil
}

// Start starts the worker
func (w *WorkerServer) Start(ctx context.Context, worker *Worker) {
	for _, task := range worker.Tasks {
		w.AddTask(task)
	}

	go func() {
		for {
			select {
			case task := <-w.taskQueue.enqueue:
				w.taskQueue.dequeue <- task
			case task := <-w.taskQueue.dequeue:
				go task.Run(ctx, w.state)
			case <-w.shutdown:
				return
			}
		}
	}()

	<-w.shutdown
	log.Info().Msg("worker server shutting down")
}

// Stop stops the worker
func (w *WorkerServer) Shutdown(ctx context.Context) error {
	close(w.taskQueue.enqueue)
	close(w.taskQueue.dequeue)

	w.shutdown <- true
	return nil
}
