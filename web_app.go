package golaze

import (
	"github.com/go-chi/chi/v5"
)

type WebAppConfig struct {
	Port string
}

type WebApp struct {
	*WebAppConfig
}

func NewWebApp(config *WebAppConfig) *WebApp {
	if config.Port == "" {
		config.Port = "8080"
	}

	return &WebApp{
		config,
	}
}

func (app *WebApp) Router() *chi.Mux {
	router := chi.NewRouter()
	return router
}
