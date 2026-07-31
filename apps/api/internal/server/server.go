package server

import (
	"log/slog"
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
)

// New builds an http.Server with timeouts from config.
func New(cfg *config.Config, handler http.Handler, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}
