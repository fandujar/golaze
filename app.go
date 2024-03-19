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
	Name        string
	Version     string
	LogLevel    zerolog.Level
	HealthCheck *HealthCheck
	WebApp      *WebApp
	Worker      *Worker
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

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-signals
		log.Info().Msgf("received signal: %v", s)
		shutdown <- true
	}()

	healthServer := &http.Server{
		Addr:           ":" + app.HealthCheck.Port,
		Handler:        app.HealthCheck.Router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ctx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	go func() {
		log.Info().Msgf("starting health check server on port %s", app.HealthCheck.Port)
		if err := healthServer.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("health check server failed")
		}

		if err := healthServer.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("health check server shutdown failed")
		}
	}()

	if app.WebApp != nil {
		webAppServer := &http.Server{
			Addr:           ":" + app.WebApp.Port,
			Handler:        app.WebApp.Router,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   5 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		go func() {
			log.Info().Msg("starting main server on port 8080")
			if err := webAppServer.ListenAndServe(); err != nil {
				log.Fatal().Err(err).Msg("main server failed")
			}

			if err := webAppServer.Shutdown(ctx); err != nil {
				log.Fatal().Err(err).Msg("main server shutdown failed")
			}
		}()
	}

	if app.Worker != nil {
		go func() {
			log.Info().Msg("starting worker")
			app.Worker.Start()
		}()
	}

	<-shutdown
	log.Info().Msg("shutting down")

	return nil
}
