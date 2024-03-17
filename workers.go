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
	lock            sync.Mutex
	Shutdown        chan bool
	State           *State
	ConcurrentTasks int
}

type Worker struct {
	*WorkerConfig
}
type State struct {
	Data map[string]interface{}
	lock sync.Mutex
}

func NewWorker(config *WorkerConfig) *Worker {
	if config.Tasks == nil {
		config.Tasks = make([]*Task, 0)
	}

	if config.Shutdown == nil {
		config.Shutdown = make(chan bool)
	}

	if config.ConcurrentTasks == 0 {
		config.ConcurrentTasks = 1
	}

	if config.State == nil {
		config.State = &State{
			Data: make(map[string]interface{}),
		}
	}

	return &Worker{
		config,
	}
}

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

// Start the worker non-blocking
func (w *Worker) Start() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		s := <-signals
		log.Info().Msgf("received signal: %v", s)
		w.Shutdown <- true
	}()

	for {
		select {
		case <-w.Shutdown:
			log.Info().Msg("worker stopped")
			return
		default:
			// Type of tasks are:
			// - tasks that run once
			// - tasks that run n times or forever
			// - tasks that can retry n times in case of failure

			if len(w.Tasks) == 0 {
				continue
			}

			// run tasks concurrently respecting the limit of concurrent tasks
			var wg sync.WaitGroup
			if len(w.Tasks) < w.ConcurrentTasks {
				w.ConcurrentTasks = len(w.Tasks)
			}
			for i := 0; i < w.ConcurrentTasks; i++ {
				task := w.Tasks[i]
				if task.Repeat == 0 {
					w.Tasks = append(w.Tasks[:i], w.Tasks[i+1:]...)
				}

				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					task.Run(w.State)
				}(i)
			}
		}
	}
}

func (w *Worker) Stop() {
	w.Shutdown <- true
}
