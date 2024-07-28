package golaze

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	Name            string
	Version         string
	LogLevel        zerolog.Level
	HealthCheck     *HealthCheck
	WebApp          *WebApp
	Worker          *Worker
	ShutdownTimeout time.Duration
}

type App struct {
	*AppConfig
}

func NewApp(config *AppConfig) *App {
	if config.LogLevel == zerolog.NoLevel {
		config.LogLevel = zerolog.InfoLevel
	}

	if config.HealthCheck == nil {
		config.HealthCheck = NewHealthCheck(
			&HealthCheckConfig{},
		)
	}

	if config.WebApp != nil && config.WebApp.Port == "" {
		config.WebApp.Port = "8080"
	}

	if config.ShutdownTimeout == 0 {
		config.ShutdownTimeout = 10 * time.Second
	}

	return &App{
		config,
	}
}

func (app *App) AddWorker(worker *Worker) {
	app.Worker = worker
}

func (app *App) Run() error {
	zerolog.SetGlobalLevel(app.LogLevel)

	shutdown := make(chan bool, 1)
	signals := make(chan os.Signal, 1)

	var healthCheckServer *http.Server
	var webAppServer *http.Server
	var workerServer *WorkerServer

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-signals
		log.Info().Msgf("received signal: %v", s)
		shutdown <- true
	}()

	healthCheckServer = &http.Server{
		Addr:           ":" + app.HealthCheck.Port,
		Handler:        app.HealthCheck.Router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Info().Msgf("starting health check server on port %s", app.HealthCheck.Port)
		if err := healthCheckServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("health check server failed")
		}
	}()

	if app.WebApp != nil {
		webAppServer = &http.Server{
			Addr:           ":" + app.WebApp.Port,
			Handler:        app.WebApp.Router,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		go func() {
			log.Info().Msgf("starting web app server on port %s", app.WebApp.Port)
			if err := webAppServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("main server failed")
			}
		}()
	}

	if app.Worker != nil {
		workerServer = NewWorkerServer()
		go func() {
			log.Info().Msg("starting worker server")
			ctx := context.Background()
			workerServer.Start(ctx, app.Worker)
		}()
	}

	<-shutdown

	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.ShutdownTimeout)
	defer cancel()

	if err := healthCheckServer.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("health check server shutdown failed")
	}

	if app.WebApp != nil && webAppServer != nil {
		if err := webAppServer.Shutdown(shutdownCtx); err != nil {
			log.Fatal().Err(err).Msg("main server shutdown failed")
		}
	}

	if app.Worker != nil && workerServer != nil {
		if err := workerServer.Shutdown(shutdownCtx); err != nil {
			log.Fatal().Err(err).Msg("worker shutdown failed")
		}
	}

	log.Info().Msg("shutdown complete")

	return nil
}
