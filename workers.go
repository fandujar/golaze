package golaze

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

type WorkerConfig struct {
	Tasks    []*Task
	Shutdown chan bool
}

type Worker struct {
	*WorkerConfig
}

type TaskConfig struct {
	Name string
	Exec func() error
	Done chan bool
}

type Task struct {
	*TaskConfig
}

func NewTask(config *TaskConfig) *Task {
	if config.Exec == nil {
		config.Exec = func() error {
			return nil
		}
	}

	if config.Done == nil {
		config.Done = make(chan bool)
	}

	return &Task{
		config,
	}
}

func (t *Task) Run() {
	go func() {
		if err := t.Exec(); err != nil {
			log.Error().Err(err).Msg("task failed")
		}

		t.Done <- true
	}()

	<-t.Done
	log.Info().Msgf("task %s completed", t.Name)
}

func NewWorker(config *WorkerConfig) *Worker {
	return &Worker{
		config,
	}
}

func (w *Worker) AddTask(task *Task) {
	w.Tasks = append(w.Tasks, task)
}

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
			for _, task := range w.Tasks {
				go func() {
					task.Run()
				}()
			}
		}
	}
}
