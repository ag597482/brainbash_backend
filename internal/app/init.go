package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"brainbash_backend/config"
	httpRouter "brainbash_backend/internal/router/http"
)

type App struct {
	httpServer *http.Server
}

func NewApp(cfg *config.AppConfig) (*App, error) {
	httpRouter.Init(cfg)

	httpServer := &http.Server{
		Addr:    ":" + cfg.StaticConfig.App.Port,
		Handler: httpRouter.Instance(),
	}

	return &App{
		httpServer: httpServer,
	}, nil
}

func (a *App) Start() error {
	errCh := make(chan error, 1)

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down server...")
	case err := <-errCh:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server stopped gracefully")
	return nil
}
