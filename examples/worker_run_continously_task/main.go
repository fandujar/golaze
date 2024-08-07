package main

import (
	"fmt"
	"time"

	"github.com/fandujar/golaze"
)

func main() {
	app := golaze.NewApp(
		&golaze.AppConfig{
			Name: "Worker Continously Run Task",
			Worker: golaze.NewWorker(
				&golaze.WorkerConfig{},
			),
		},
	)

	task := golaze.NewTask(
		&golaze.TaskConfig{
			Name:        "task run forever",
			Repeat:      -1,
			RepeatDelay: 10 * time.Second,
			Timeout:     15 * time.Second,
			Exec: func(state *golaze.State, cancel chan bool) error {
				if state.Get("counter") == nil {
					state.Set("counter", 0)
				}

				state.Set("counter", state.Get("counter").(int)+1)
				fmt.Printf("running task example: %d\n", state.Get("counter"))
				time.Sleep(2 * time.Second)
				return nil
			},
		})

	app.Worker.Tasks = append(app.Worker.Tasks, task)
	app.Run()
}
