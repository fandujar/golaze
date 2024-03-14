package golaze

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	LogLevel    *zerolog.Level
	HealthCheck *HealthCheck
}

type App struct {
	*AppConfig
}

func NewApp(config *AppConfig) *App {
	if config.LogLevel == nil {
		level := zerolog.InfoLevel
		config.LogLevel = &level
	}

	if config.HealthCheck == nil {
		config.HealthCheck = NewHealthCheck(
			&HealthCheckConfig{},
		)
	}

	return &App{
		config,
	}
}

func (app *App) Run() error {
	zerolog.SetGlobalLevel(*app.LogLevel)

	shutdown := make(chan bool, 1)
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-signals
		log.Info().Msgf("received signal: %v", s)
		shutdown <- true
	}()

	<-shutdown
	log.Info().Msg("shutting down")

	return nil
}
