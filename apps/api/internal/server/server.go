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
		Handler:           handler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: headerTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}
