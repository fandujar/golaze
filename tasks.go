package golaze

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type TaskConfig struct {
	Name          string
	Exec          func(state *State, cancel chan bool) error
	MaxRetries    int
	RetryInterval time.Duration
	Repeat        int // -1 for infinite, 0 for no repeat, > 0 for n times
	RepeatDelay   time.Duration
	Timeout       time.Duration
	RunHistory    []time.Time

	Cancel chan bool
	Done   chan bool
}

type Task struct {
	*TaskConfig
}

func NewTask(config *TaskConfig) *Task {
	if config.Cancel == nil {
		config.Cancel = make(chan bool)
	}

	if config.Exec == nil {
		config.Exec = func(state *State, cancel chan bool) error {
			return nil
		}
	}

	if config.Done == nil {
		config.Done = make(chan bool)
	}

	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	if config.RetryInterval == 0 {
		config.RetryInterval = 5 * time.Second
	}

	if config.RepeatDelay == 0 {
		config.RepeatDelay = 1 * time.Second
	}

	if config.Repeat == 0 {
		config.Repeat = 1
	}

	return &Task{
		config,
	}
}

func (t *Task) IsRunning() bool {
	// Check if the task is currently running
	// by checking the length of the Done channel
	select {
	case <-t.Done:
		return false
	default:
		return true
	}
}

func (t *Task) Run(state *State, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		t.RunHistory = append(t.RunHistory, time.Now())
		log.Info().Msgf("task %s started", t.Name)
		if err := t.Exec(state, t.Cancel); err != nil {
			log.Error().Err(err).Msgf("task %s failed", t.Name)
		}

		t.Done <- true
		wg.Done()
	}()

	select {
	case <-t.Cancel:
		log.Info().Msgf("task %s cancelled", t.Name)
		wg.Done()
		t.Done <- true

	case <-t.Done:
		log.Info().Msgf("task %s completed", t.Name)
		wg.Done()

	case <-time.After(t.Timeout):
		log.Error().Msgf("task %s timed out", t.Name)
		wg.Done()
		t.Done <- true
	}
}
