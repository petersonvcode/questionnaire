package server

import (
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
)

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Uncomment to debug CORS preflight requests
		// if r.Method == http.MethodOptions {
		// 	requestedHeaders := r.Header.Get("Access-Control-Request-Headers")
		// 	slog.Debug("CORS Preflight - Requested Headers:", "requestedHeaders", requestedHeaders)
		// }

		origin := r.Header.Get("Origin")
		if isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Token, Cache-Control, Credentials")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else if origin != "" {
			slog.Warn("CORS - Origin not allowed", "origin", origin, "method", r.Method, "path", r.URL.Path)
		}

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	allowedOrigins := getAllowedOrigin()
	if len(allowedOrigins) == 0 {
		panic("CORS_ALLOWED_ORIGINS is not set")
	}
	return allowedOrigins[0] == "*" || slices.Contains(allowedOrigins, origin)
}

var splitOrigins []string

func getAllowedOrigin() []string {
	if splitOrigins == nil {
		splitOrigins = strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	}
	return splitOrigins
}
