package golaze

import (
	"github.com/go-chi/chi/v5"
)

type WebAppConfig struct {
	Port   string
	Router *chi.Mux
}

type WebApp struct {
	*WebAppConfig
}

func NewWebApp(config *WebAppConfig) *WebApp {
	if config.Port == "" {
		config.Port = "8080"
	}

	if config.Router == nil {
		config.Router = NewRouter()
	}

	return &WebApp{
		config,
	}
}
