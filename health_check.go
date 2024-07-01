package golaze

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type HealthCheckConfig struct {
	Port           string
	LivenessHooks  []func() error
	ReadinessHooks []func() error
	Router         *chi.Mux
}

type HealthCheck struct {
	*HealthCheckConfig
}

func LivenessHandler(hooks ...func() error) func(w http.ResponseWriter, r *http.Request) {
	errors := make(chan error, len(hooks))
	// Run all the hooks concurrently
	for _, hook := range hooks {
		go func(hook func() error) {
			errors <- hook()
		}(hook)
	}

	// Wait for all the hooks to finish
	for i := 0; i < len(hooks); i++ {
		if err := <-errors; err != nil {
			return func(w http.ResponseWriter, r *http.Request) {
				JSONError(w, "error", http.StatusInternalServerError)
			}
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(w, map[string]string{"status": "ok"}, http.StatusOK)
	}
}

func ReadinessHandler(hooks ...func() error) func(w http.ResponseWriter, r *http.Request) {
	errors := make(chan error, len(hooks))
	// Run all the hooks concurrently
	for _, hook := range hooks {
		go func(hook func() error) {
			errors <- hook()
		}(hook)
	}

	// Wait for all the hooks to finish
	for i := 0; i < len(hooks); i++ {
		if err := <-errors; err != nil {
			return func(w http.ResponseWriter, r *http.Request) {
				JSONError(w, "error", http.StatusInternalServerError)
			}
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(w, map[string]string{"status": "ok"}, http.StatusOK)
	}
}

func NewHealthCheck(config *HealthCheckConfig) *HealthCheck {
	if config.Port == "" {
		port := os.Getenv("HEALTHCHECK_PORT")
		if port == "" {
			config.Port = "8081"
		} else {
			config.Port = port
		}
	}

	if config.Router == nil {
		r := NewRouter()
		r.Get("/liveness", LivenessHandler(config.LivenessHooks...))
		r.Get("/readiness", ReadinessHandler(config.ReadinessHooks...))

		config.Router = r
	}

	return &HealthCheck{
		config,
	}
}
