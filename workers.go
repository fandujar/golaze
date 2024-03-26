package golaze

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type WorkerConfig struct {
	Tasks           []*Task
	Runners         []*Runner
	Shutdown        chan bool
	State           *State
	ConcurrentTasks int
	lock            sync.Mutex
}

type Worker struct {
	*WorkerConfig
	taskQueue *TaskQueue
}

type Runner struct {
	ID int
}

func NewRunner(id int) *Runner {
	return &Runner{
		ID: id,
	}
}

func (r *Runner) Run(taskQueue chan Task, ctx context.Context, state *State) {
	log.Debug().Msgf("runner %d started", r.ID)
	for {
		select {
		case task := <-taskQueue:
			if task == (Task{}) {
				continue
			}
			log.Debug().Msgf("runner %d running task %s", r.ID, task.Name)
			task.Run(ctx, state)
		case <-ctx.Done():
			log.Debug().Msgf("runner %d stopped", r.ID)
			return
		}
	}
}

func NewWorker(config *WorkerConfig) *Worker {
	if config.Tasks == nil {
		config.Tasks = make([]*Task, 0)
	}

	if config.Shutdown == nil {
		config.Shutdown = make(chan bool)
	}

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

	return &Worker{
		config,
		queue,
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

	w.taskQueue.Start()
	w.Runners = make([]*Runner, w.ConcurrentTasks)

	for i := 0; i < w.ConcurrentTasks; i++ {
		w.Runners[i] = NewRunner(i)
		ctx := w.State.Context()
		go w.Runners[i].Run(w.taskQueue.dequeue, ctx, w.State)
	}

	for {
		select {
		case <-w.Shutdown:
			log.Info().Msg("worker stopped")
			return
		default:
			for _, task := range w.Tasks {
				go func(t *Task) {
					for {
						select {
						case <-w.Shutdown:
							return
						case <-t.Done:
							<-time.After(t.RetryInterval)
							if t.Repeat >= 1 {
								t.Repeat--
							}

							if t.Repeat == 0 {
								return
							}

							w.taskQueue.enqueue <- *t
						default:
							if t.Done == nil {
								t.Done = make(chan bool)
								t.Done <- false
							}
						}
					}
				}(task)
			}
		}
	}
}

func (w *Worker) Stop() {
	w.Shutdown <- true
}
