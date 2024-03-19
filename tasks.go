package golaze

import (
	"time"

	"github.com/rs/zerolog/log"
)

type TaskConfig struct {
	Name          string
	Exec          func(state *State) error
	Done          chan bool
	MaxRetries    int
	RetryInterval time.Duration
	Repeat        int // -1 for infinite, 0 for no repeat, > 0 for n times
	RepeatDelay   time.Duration
	Timeout       time.Duration
	RunHistory    []time.Time
}

type Task struct {
	*TaskConfig
}

func NewTask(config *TaskConfig) *Task {
	if config.Exec == nil {
		config.Exec = func(state *State) error {
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

func (t *Task) Run(state *State) {
	go func() {
		t.RunHistory = append(t.RunHistory, time.Now())
		if err := t.Exec(state); err != nil {
			log.Error().Err(err).Msgf("task %s failed", t.Name)
		}

		t.Done <- true
	}()

	select {
	case <-t.Done:
		log.Info().Msgf("task %s completed", t.Name)
	case <-time.After(t.Timeout):
		log.Error().Msgf("task %s timed out", t.Name)
		t.Done <- true
	}
}
