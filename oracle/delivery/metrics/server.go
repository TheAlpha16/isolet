package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/TheAlpha16/isolet/oracle/utils"
	"github.com/TheAlpha16/isolet/oracle/utils/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	logger *logger.StandardLogger
}

func New() *Server {
	config := utils.GetConfig()
	mux := http.NewServeMux()
	mux.Handle(config.Metrics.Path, promhttp.Handler())

	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", config.Metrics.Port),
			Handler: mux,
		},
		logger: logger.GetAppLogger(),
	}
}

func (s *Server) Start() {
	config := utils.GetConfig()
	s.logger.Info("Starting metrics server", zap.Int("port", config.Metrics.Port), zap.String("path", config.Metrics.Path))

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("Failed to start metrics server", zap.Error(err))
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down metrics server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Metrics server shutdown failed", zap.Error(err))
		return err
	}

	s.logger.Info("Metrics server shutdown successful")
	return nil
}
