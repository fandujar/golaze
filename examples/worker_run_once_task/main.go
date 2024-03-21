package main

import (
	"fmt"
	"time"

	"github.com/fandujar/golaze"
)

func main() {
	worker := golaze.NewWorker(
		&golaze.WorkerConfig{},
	)

	task := golaze.NewTask(
		&golaze.TaskConfig{
			Name:    "example 1",
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
			Name:    "example 2",
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

	task3 := golaze.NewTask(
		&golaze.TaskConfig{
			Name:    "example 3",
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

	worker.AddTask(task)
	go worker.Start()

	err := worker.AddTask(task2)
	if err != nil {
		fmt.Println(err)
	}

	err = worker.AddTask(task3)
	if err != nil {
		fmt.Println(err)
	}

	// wait for the worker to finish
	<-worker.Shutdown
	fmt.Println("worker finished")
}
