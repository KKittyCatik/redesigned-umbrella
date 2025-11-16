package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/KKittyCatik/redesigned-umbrella/internal/infrastructure/config"
)

func StartServer(ctx context.Context, logger *slog.Logger, serverCfg config.ServerConfig, handler http.Handler) {
	server := &http.Server{
		Addr:         ":" + serverCfg.Port,
		Handler:      handler,
		ReadTimeout:  serverCfg.ReadTimeout,
		WriteTimeout: serverCfg.WriteTimeout,
		IdleTimeout:  serverCfg.IdleTimeout,
	}

	go func() {
		logger.Info("Server starting",
			"addr", server.Addr,
			"readTimeout", serverCfg.ReadTimeout,
			"writeTimeout", serverCfg.WriteTimeout,
			"idleTimeout", serverCfg.IdleTimeout,
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown error", "error", err)
	} else {
		logger.Info("Server stopped gracefully")
	}
}
