package golaze

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
)

type WorkerConfig struct {
	Tasks           []*Task
	Shutdown        chan bool
	State           *State
	ConcurrentTasks int
	lock            sync.Mutex
}

type Worker struct {
	*WorkerConfig
	taskQueue *TaskQueue
}

type TaskQueue struct {
	tasks   []Task
	enqueue chan Task
	dequeue chan Task
}

type State struct {
	Data map[string]interface{}
	lock sync.Mutex
}

func (s *State) Set(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Data[key] = value
}

func (s *State) Get(key string) interface{} {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Data[key]
}

func (s *State) Delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.Data, key)
}

func (s *State) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Data = make(map[string]interface{})
}

func (s *State) Context() context.Context {
	return context.WithValue(context.Background(), "state", s)
}

func (tq *TaskQueue) Start() {
	go func() {
		for {
			if len(tq.tasks) == 0 {
				continue
			}

			select {
			case task := <-tq.enqueue:
				log.Debug().Msgf("enqueuing task %s", task.Name)
				tq.tasks = append(tq.tasks, task)
			case tq.dequeue <- tq.tasks[0]:
				log.Debug().Msgf("dequeuing task %s", tq.tasks[0].Name)
				tq.tasks = tq.tasks[1:]
			}
		}
	}()
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

	for i := 0; i < w.ConcurrentTasks; i++ {
		go func() {
			log.Info().Msg("runner started")
			for {
				select {
				case <-w.Shutdown:
					return
				default:
					task := <-w.taskQueue.dequeue
					task.Run(w.State)
				}
			}
		}()
	}

	for {
		select {
		case <-w.Shutdown:
			log.Info().Msg("worker stopped")
			return
		default:
			for _, task := range w.Tasks {
				w.taskQueue.enqueue <- *task
				// if task.Repeat == -1 {
				// 	go func() {
				// 		time.After(task.RepeatDelay)
				// 		w.taskQueue.enqueue <- *task
				// 	}()
				// }
				// if task.Repeat >= 1 {
				// 	task.Repeat--
				// 	go func() {
				// 		time.After(task.RepeatDelay)
				// 		w.taskQueue.enqueue <- *task
				// 	}()
				// }
			}
		}
	}
}

func (w *Worker) Stop() {
	w.Shutdown <- true
}
