package main

import (
	"net/http"

	"github.com/fandujar/golaze"
	"github.com/rs/zerolog"
)

func main() {
	router := golaze.NewRouter()
	router.Use(golaze.LogMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		golaze.JSONResponse(w, map[string]string{"message": "Hello, World!"}, http.StatusOK)
	})

	webapp := golaze.NewWebApp(
		&golaze.WebAppConfig{
			Port:   "8080",
			Router: router,
		},
	)

	app := golaze.NewApp(
		&golaze.AppConfig{
			LogLevel: zerolog.DebugLevel,
			WebApp:   webapp,
		},
	)

	app.Run()

}
