package golaze

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type AuthToken struct {
	Token string
}

func NewAuthMiddleware(getAuthTokens func() ([]*AuthToken, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authTokens, err := getAuthTokens()
			if err != nil {
				JSONError(w, "error getting auth tokens", http.StatusInternalServerError)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				JSONError(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			authorized := false
			for _, t := range authTokens {
				if t.Token == token {
					authorized = true
					break
				}
			}

			if !authorized {
				JSONError(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		rctx := chi.RouteContext(r.Context())

		host := r.Header.Get("Host")
		url := r.URL.String()

		log.Info().Msgf("request: %v | method: %v | host: %v | status: %v | route pattern: %v", url, r.Method, host, ww.Status(), rctx.RoutePatterns)
	})
}
