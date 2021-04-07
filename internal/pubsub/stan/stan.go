package stan

import (
	"context"
	"sync"
	"time"

	"ptcg_trader/internal/config"

	"github.com/cenk/backoff"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

var (
	waitClose sync.WaitGroup
)

// NewStanConn ...
func NewStanConn(cfg config.StanConfig) (stan.Conn, error) {
	waitClose.Add(1)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 3 * time.Minute
	var stanConn stan.Conn
	err := backoff.Retry(func() error {
		var err error

		// nats connection
		natsConn, err := nats.Connect(cfg.Address,
			nats.ClosedHandler(func(_ *nats.Conn) {
				waitClose.Done()
			}),
		)
		if err != nil {
			log.Error().Msgf("fail to connect to nats: %s, addr: %v", err.Error(), cfg.Address)
			return err
		}

		if cfg.ClientID == "" {
			cfg.ClientID = xid.New().String()
		} else {
			cfg.ClientID += "_" + xid.New().String()
		}

		// nats streaming connection
		stanConn, err = stan.Connect(cfg.ClusterID, cfg.ClientID,
			stan.NatsConn(natsConn),
			stan.Pings(10, 5),
			stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
				log.Error().Msgf("Connection lost, reason: %v", reason)
			}),
		)
		if err != nil {
			log.Error().Msgf("fail to connect to stan: %s, clusterID: %v, clientID: %v", err.Error(), cfg.ClusterID, cfg.ClientID)
			return err
		}

		return nil
	}, bo)
	if err != nil {
		return nil, err
	}

	return stanConn, nil
}

// NewStanConnWithFX ...
func NewStanConnWithFX(lc fx.Lifecycle, cfg config.StanConfig) (stan.Conn, error) {
	stanConn, err := NewStanConn(cfg)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			err := stanConn.NatsConn().Drain()
			if err != nil {
				return err
			}
			waitClose.Wait()
			return nil
		},
	})

	return stanConn, nil
}
