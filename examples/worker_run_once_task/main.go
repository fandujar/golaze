package main

import (
	"fmt"
	"time"

	"github.com/fandujar/golaze"
)

func main() {
	app := golaze.NewApp(
		&golaze.AppConfig{
			Name: "Worker Run Once Task",
			Worker: golaze.NewWorker(
				&golaze.WorkerConfig{},
			),
		},
	)

	task := golaze.NewTask(
		&golaze.TaskConfig{
			Name:    "task 1 - complete",
			Timeout: 3 * time.Second,
			Exec: func(state *golaze.State, cancel chan bool) error {
				if state.Get("example") == nil {
					state.Set("example", 0)
				}

				state.Set("example", state.Get("example").(int)+1)
				fmt.Printf("running task example: %d\n", state.Data["example"])

				return nil
			},
		})

	task2 := golaze.NewTask(
		&golaze.TaskConfig{
			Name:    "task 2 - timeout",
			Timeout: 3 * time.Second,
			Exec: func(state *golaze.State, cancel chan bool) error {
				if state.Get("example") == nil {
					state.Set("example", 0)
				}

				state.Set("example", state.Get("example").(int)+1)
				fmt.Printf("running task example: %d\n", state.Data["example"])

				time.Sleep(5 * time.Second)
				return nil
			},
		})

	task3 := golaze.NewTask(
		&golaze.TaskConfig{
			Name:    "task 3 - cancel",
			Timeout: 3 * time.Second,
			Exec: func(state *golaze.State, cancel chan bool) error {
				if state.Get("example") == nil {
					state.Set("example", 0)
				}

				state.Set("example", state.Get("example").(int)+1)
				fmt.Printf("running task example: %d\n", state.Data["example"])

				cancel <- true

				return nil
			},
		})

	app.Worker.Tasks = append(app.Worker.Tasks, task)
	app.Worker.Tasks = append(app.Worker.Tasks, task2)
	app.Worker.Tasks = append(app.Worker.Tasks, task3)

	app.Run()
}
