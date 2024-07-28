package main

import (
	"fmt"
	"time"

	"github.com/fandujar/golaze"
	"github.com/rs/zerolog/log"
)

func main() {
	app := golaze.NewApp(
		&golaze.AppConfig{
			Name: "Worker Run Retry Task",
			Worker: golaze.NewWorker(
				&golaze.WorkerConfig{},
			),
		},
	)

	task := golaze.NewTask(
		&golaze.TaskConfig{
			Name:          "task retry n times",
			MaxRetries:    3,
			RetryInterval: 5 * time.Second,
			Timeout:       15 * time.Second,
			Exec: func(state *golaze.State, cancel chan bool) error {
				if state.Get("counter") == nil {
					state.Set("counter", 0)
				}

				log.Info().Msgf("running task example: %d", state.Get("counter"))
				if state.Get("counter").(int) < 3 {
					state.Set("counter", state.Get("counter").(int)+1)
					return fmt.Errorf("error on task")
				}
				return nil
			},
		})

	app.Worker.Tasks = append(app.Worker.Tasks, task)
	app.Run()
}
