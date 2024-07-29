package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fandujar/golaze"
)

func main() {
	app := golaze.NewApp(
		&golaze.AppConfig{
			Name: "Web App",
		},
	)

	app.WebApp = golaze.NewWebApp(
		&golaze.WebAppConfig{},
	)

	app.Worker = golaze.NewWorker(
		&golaze.WorkerConfig{
			EventBus: app.EventBus,
		},
	)

	app.WebApp.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		task := golaze.NewTask(
			&golaze.TaskConfig{
				Name:    "task once run thru http",
				Timeout: 3 * time.Second,
				Exec: func(state *golaze.State, cancel chan bool) error {
					fmt.Println("running task example")

					return nil
				},
			})

		event := &golaze.Event{
			Data: task,
		}

		app.EventBus.Publish("task", event)
	})

	app.Run()

}
