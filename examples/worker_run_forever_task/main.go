package main

import (
	"fmt"

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
			Exec: func(state *golaze.State) error {
				if state.Get("counter") == nil {
					state.Set("counter", 0)
				}

				state.Set("counter", state.Get("counter").(int)+1)
				fmt.Printf("running task example: %d\n", state.Get("counter"))
				// time.Sleep(10 * time.Second)
				return nil
			},
		})

	worker.AddTask(task)
	go worker.Start()

	// wait for the worker to finish
	<-worker.Shutdown
	fmt.Println("worker finished")
}
