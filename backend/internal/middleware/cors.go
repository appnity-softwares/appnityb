package middleware

import (
	"net/http"

	"github.com/appnity/backend/internal/config"
)

type CORSMiddleware struct {
	cfg config.Config
}

func NewCORSMiddleware(cfg config.Config) *CORSMiddleware {
	return &CORSMiddleware{cfg: cfg}
}

func (m *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		isAllowed := false
		isWildcard := false

		if origin != "" {
			for _, allowed := range m.cfg.CORSOrigins {
				if allowed == "*" {
					isWildcard = true
					isAllowed = true
					break
				}
				if allowed == origin {
					isAllowed = true
					break
				}
			}
		}

		if isAllowed {
			if isWildcard {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}

		allowedMethods := "GET, POST, PUT, PATCH, DELETE, OPTIONS"
		allowedHeaders := "Content-Type, Authorization, X-Requested-With"

		w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
