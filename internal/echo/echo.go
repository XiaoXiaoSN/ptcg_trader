package echo

import (
	"context"
	"net/http"
	"time"

	"ptcg_trader/internal/config"
	"ptcg_trader/internal/echo/middleware"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

// NewEchoEngine create new echo engine for handler to register
func NewEchoEngine(cfg config.HTTPConfig) *echo.Echo {
	echo.NotFoundHandler = notFoundHandler
	echo.MethodNotAllowedHandler = notFoundHandler

	e := echo.New()

	if cfg.Mode == "debug" {
		e.Debug = true
		e.HideBanner = false
		e.HidePort = false
	} else {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true
	}
	e.HTTPErrorHandler = httpErrorHandler

	// all of API have to take a traceID in ctx
	e.Pre(middleware.NewTraceIDMiddleware())

	// allow cors
	e.Use(middleware.CORSConfig)

	// record response errors form echo service
	e.Use(middleware.ErrMiddleware)

	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// register default route
	RegisterDefaultRoute(e)

	return e
}

// StartEcho create new echo engine for handler to register
func StartEcho(s config.HTTPConfig, lc fx.Lifecycle) *echo.Echo {
	e := NewEchoEngine(s)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info().Msgf("Starting echo server, listen on %s", s.Address)
			go func() {
				err := e.Start(s.Address)
				if err != nil {
					log.Error().Msgf("Error echo server, err: %s", err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Stopping echo HTTP server.")
			return e.Shutdown(ctx)
		},
	})

	return e
}

// RegisterDefaultRoute provide default handler
func RegisterDefaultRoute(e *echo.Echo) {
	e.GET("/ping", func(c echo.Context) error {
		time.Sleep(10 * time.Second)
		return c.String(http.StatusOK, "pong!!!")
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!!!")
	})
}
