package obs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// ServerOptions carries the raw observability.metrics config values (obs does
// not import the config package; composition maps Config fields onto these).
type ServerOptions struct {
	Enabled   bool
	Addr      string
	AuthToken string
}

// Server is the optional dedicated metrics listener (R1 D-001 s1, R2 D-001
// s3). Disabled means Start is a no-op and nothing listens - the documented
// no-collector default. Enabled bind failures fail startup (explicit config is
// a hard requirement); a runtime Serve failure only logs loudly because the
// bypass face must never take the API process down.
type Server struct {
	enabled   bool
	addr      string
	handler   http.Handler
	logger    *slog.Logger
	srv       *http.Server
	boundAddr string
}

// NewServer assembles the listener. A nil observer disables the server even
// when options claim enabled.
func NewServer(opts ServerOptions, observer *Observer, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	s := &Server{enabled: opts.Enabled && observer != nil, addr: opts.Addr, logger: logger}
	if s.enabled {
		mux := http.NewServeMux()
		mux.Handle("GET "+MetricsPath, observer.Handler(opts.AuthToken))
		s.handler = mux
	}
	return s
}

// Enabled reports whether the listener will actually serve.
func (s *Server) Enabled() bool { return s != nil && s.enabled }

// Addr returns the bound address after a successful Start ("" otherwise).
func (s *Server) Addr() string {
	if s == nil {
		return ""
	}
	return s.boundAddr
}

// Start binds and serves the exposition face. Disabled servers return nil.
func (s *Server) Start(ctx context.Context) error {
	if !s.Enabled() {
		return nil
	}
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		// Fail-closed (R2 D-001 s3): an explicitly enabled surface that cannot
		// bind is a configuration error, not something to degrade silently.
		return fmt.Errorf("observability: metrics listen %s: %w", s.addr, err)
	}
	s.boundAddr = ln.Addr().String()
	s.srv = &http.Server{
		Handler:           s.handler,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	go func() {
		if err := s.srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("observability metrics server failed", "err", err)
		}
	}()
	s.logger.Info("observability metrics listening", "addr", s.boundAddr, "path", MetricsPath)
	return nil
}

// Stop shuts the listener down; disabled or never-started servers return nil.
func (s *Server) Stop(ctx context.Context) error {
	if s == nil || s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}
