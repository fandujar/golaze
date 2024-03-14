package golaze

import "net/http"

type HealthCheckConfig struct {
	Port           string
	LivenessHooks  []func() error
	ReadinessHooks []func() error
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
		config.Port = "8081"
	}

	return &HealthCheck{
		config,
	}
}

func (hc *HealthCheck) Router() http.Handler {
	router := chi.NewRouter()

	router.Get("/liveness", LivenessHandler(hc.LivenessHooks...))
	router.Get("/readiness", ReadinessHandler(hc.ReadinessHooks...))

	return router
}
