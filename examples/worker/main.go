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
			Name: "example",
			Exec: func() error {
				fmt.Println("hello world")
				return nil
			},
		})

	worker.Tasks = append(worker.Tasks, task)

	worker.Start()
}
