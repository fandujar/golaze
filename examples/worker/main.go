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
			Name: "example 1",
			Exec: func(state *golaze.State) error {
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

	// add the same task again after 5 seconds
	time.Sleep(5 * time.Second)
	err := worker.AddTask(task)
	if err != nil {
		fmt.Println(err)
	}

	// add the same task again after 5 seconds
	time.Sleep(5 * time.Second)
	task.Repeat = -1
	err = worker.AddTask(task)
	if err != nil {
		fmt.Println(err)
	}

	// wait for the worker to finish
	<-worker.Shutdown
	fmt.Println("worker finished")
}
