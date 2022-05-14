package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	log "github.com/rs/zerolog/log"
	fx "go.uber.org/fx"
)

func cmdRecover() {
	if r := recover(); r != nil {
		var msg string
		for i := 2; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			msg = msg + fmt.Sprintf("%s:%d\n", file, line)
		}
		log.Error().Msgf("%s\n↧↧↧↧↧↧ PANIC ↧↧↧↧↧↧\n%s↥↥↥↥↥↥ PANIC ↥↥↥↥↥↥", r, msg)
	}
}

func startGraceApp(name string, app *fx.App) error {
	if err := app.Start(context.Background()); err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info().Msgf("main: shutting down %s...", name)

	stopCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return nil
}
