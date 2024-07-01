package main

import (
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

	app.WebApp.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.Write([]byte("Hello, World!"))
	})

	app.Run()

}
