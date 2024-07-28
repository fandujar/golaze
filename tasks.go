package golaze

import (
	"context"
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

	lock sync.Mutex
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

func (t *Task) Run(ctx context.Context, state *State) {
	taskError := make(chan error)
	go func() {
		go func() {
			t.lock.Lock()
			t.RunHistory = append(t.RunHistory, time.Now())
			t.lock.Unlock()

			log.Info().Msgf("task %s started", t.Name)
			err := t.Exec(state, t.Cancel)
			taskError <- err
		}()

		select {
		case <-ctx.Done():
			log.Info().Msgf("task %s stopped", t.Name)
		case <-t.Cancel:
			log.Info().Msgf("task %s cancelled", t.Name)
		case err := <-taskError:
			if err != nil {
				log.Error().Err(err).Msgf("task %s failed", t.Name)
			} else {
				log.Info().Msgf("task %s completed", t.Name)
			}
		case <-time.After(t.Timeout):
			log.Error().Msgf("task %s timed out", t.Name)
		}

		t.Done <- true
	}()

	<-t.Done
	if taskError != nil && t.MaxRetries > 0 {
		t.MaxRetries--
		time.Sleep(t.RetryInterval)
		t.Run(ctx, state)
	}

	if t.Repeat > 0 {
		t.Repeat--
		time.Sleep(t.RepeatDelay)
		t.Run(ctx, state)
	}

	if t.Repeat == 0 {
		close(t.Done)
	}

	if t.Repeat == -1 {
		time.Sleep(t.RepeatDelay)
		t.Run(ctx, state)
	}
}
