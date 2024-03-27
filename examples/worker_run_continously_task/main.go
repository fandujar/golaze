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
			Name:        "task run forever",
			Repeat:      -1,
			RepeatDelay: 10,
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

	worker.AddTask(task)
	go worker.Start()

	// wait for the worker to finish
	worker.WaitShutdown()
}
