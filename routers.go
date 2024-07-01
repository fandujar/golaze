package golaze

import "github.com/go-chi/chi/v5"

func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
