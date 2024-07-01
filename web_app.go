package golaze

import (
	"os"

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
		port := os.Getenv("PORT")
		if port == "" {
			config.Port = "8080"
		} else {
			config.Port = port
		}
	}

	if config.Router == nil {
		config.Router = NewRouter()
	}

	return &WebApp{
		config,
	}
}
