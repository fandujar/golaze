package golaze

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

type WorkerConfig struct {
	Task     *[]Task
	Shutdown chan bool
}

type Worker struct {
	*WorkerConfig
}

type Task struct {
	Name string
	Exec func() error
	Done chan bool
}

func NewTask(name string, exec func() error) *Task {
	if exec == nil {
		exec = func() error {
			return nil
		}
	}

	return &Task{
		Name: name,
		Exec: exec,
		Done: make(chan bool, 1),
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
			for _, task := range *w.Task {
				go func() {
					task.Run()
				}()
			}
		}
	}
}
