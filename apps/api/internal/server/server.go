package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
)

// New builds an http.Server with timeouts from config.
//
// ReadHeaderTimeout is always set (defaulting to ReadTimeout, floored at 5s) so
// a client that dribbles request headers cannot hold a connection open past the
// configured budget (Slowloris-class mitigation; GOAL-003).
func New(cfg *config.Config, handler http.Handler, logger *slog.Logger) *http.Server {
	headerTimeout := cfg.ReadTimeout
	if headerTimeout <= 0 {
		headerTimeout = 5 * time.Second
	}
	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           WrapSecurity(cfg, handler),
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: headerTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}

// WrapSecurity adds nosniff and optional CORS (W15-F05). Empty CORS list
// leaves cross-origin preflight unchanged (same-origin Nginx default).
func WrapSecurity(cfg *config.Config, next http.Handler) http.Handler {
	allow := map[string]struct{}{}
	for _, origin := range cfg.HTTPCORSOrigins {
		allow[origin] = struct{}{}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		origin := r.Header.Get("Origin")
		if _, ok := allow[origin]; origin != "" && ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept-Language, X-Refresh-Token")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
