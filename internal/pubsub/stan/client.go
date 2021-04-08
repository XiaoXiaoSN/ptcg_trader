package stan

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"ptcg_trader/internal/config"
	"ptcg_trader/internal/ctxutil"
	"ptcg_trader/internal/errors"

	"github.com/cenk/backoff"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

// MockClient ...
// TODO: should use interface here !
type MockClient interface {}

// NewClient ...
func NewClient(cfg config.StanConfig) (*Client, error) {
	var client Client
	sc, err := NewStanConn(cfg)
	if err != nil {
		return nil, err
	}
	client.stanConn = sc
	client.cfg = cfg

	client.stanConn.NatsConn().SetReconnectHandler(client.reconnect)

	return &client, nil
}

// NewClientWithFx ...
func NewClientWithFx(lc fx.Lifecycle, cfg config.StanConfig) (*Client, error) {
	sc, err := NewStanConnWithFX(lc, cfg)
	if err != nil {
		return nil, err
	}

	var c Client
	c.stanConn = sc
	c.cfg = cfg

	c.stanConn.NatsConn().SetReconnectHandler(c.reconnect)

	return &c, nil
}

// Client is a stan subscriber, publisher
type Client struct {
	stanConn stan.Conn
	cfg      config.StanConfig

	Channels []Channel
}

// Pub publish message
func (c *Client) Pub(ctx context.Context, topic string, data interface{}) error {
	var msgData MsgData

	b, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(errors.ErrBadRequest, "fail to marshal pub data, err: %s", err.Error())
	}
	msgData.Data = b
	msgData.TraceID = ctxutil.TraceIDFromCtx(ctx)

	b, err = json.Marshal(msgData)
	if err != nil {
		return errors.Wrapf(errors.ErrBadRequest, "fail to marshal internal MsgData, err: %s", err.Error())
	}

	if err := c.stanConn.Publish(topic, b); err != nil {
		return errors.Wrapf(errors.ErrInternalError, "fail to publish to stan, err: %s", err.Error())
	}

	return nil
}

// RegisterChannel ...
func (c *Client) RegisterChannel(ctx context.Context, channels []Channel) error {
	c.Channels = channels

	for i := range channels {
		log.Info().Msgf("Register channel: %s", channels[i].ChannelName)

		if c.cfg.DurableName != "" {
			channels[i].Options = append(channels[i].Options, stan.DurableName(c.cfg.DurableName))
		}

		name, group, handler := channels[i].ChannelName, channels[i].GroupName, channels[i].Handler
		_, err := c.stanConn.QueueSubscribe(name, group, func(msg *stan.Msg) {
			defer recoverLog()
			logger := log.Ctx(ctx).With().Str("endpoint", name).Logger()

			var msgData MsgData
			err := json.Unmarshal(msg.Data, &msgData)
			if err != nil {
				logger.Error().Msgf("Fail to unmarshal to internal msgData: %s", msgData)
				return
			}

			logger = logger.With().
				Str("trace_id", msgData.TraceID).
				Logger()

			if c.cfg.Debug {
				logger.Debug().
					Str("stan_channel", name).
					Str("stan_group", group).
					Str("stan_req", string(msgData.Data)).
					Msg("access log")
			}
			_ctx := logger.WithContext(ctx)
			_ctx = context.WithValue(_ctx, ctxutil.CtxKeyTraceID, msgData.TraceID)

			err = handler(_ctx, msgData.Data)
			if err != nil {
				logger.Error().Msgf("channel: %s, error: %+v", name, err)
			}

		}, channels[i].Options...)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisterNatsChannel ...
func (c *Client) RegisterNatsChannel(ctx context.Context, channels []Channel) error {
	for i := range channels {
		log.Info().Msgf("Register channel: %s", channels[i].ChannelName)

		if c.cfg.DurableName != "" {
			channels[i].Options = append(channels[i].Options, stan.DurableName(c.cfg.DurableName))
		}
		name, group, handler := channels[i].ChannelName, channels[i].GroupName, channels[i].Handler
		_, err := c.stanConn.NatsConn().QueueSubscribe(name, group, func(msg *nats.Msg) {
			defer recoverLog()
			logger := log.Ctx(ctx).With().Str("endpoint", name).Logger()

			var msgData MsgData
			err := json.Unmarshal(msg.Data, &msgData)
			if err != nil {
				logger.Error().Msgf("Fail to unmarshal to internal msgData: %s", msgData)
				return
			}

			logger = logger.With().
				Str("trace_id", msgData.TraceID).
				Logger()
			_ctx := logger.WithContext(ctx)
			_ctx = context.WithValue(_ctx, ctxutil.CtxKeyTraceID, msgData.TraceID)

			err = handler(_ctx, msgData.Data)
			if err != nil {
				logger.Error().Msgf("channel: %s, error: %+v", name, err)
			}

		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) reconnect(conn *nats.Conn) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 3 * time.Minute
	err := backoff.Retry(func() error {
		var err error
		c.stanConn, err = stan.Connect(c.cfg.ClusterID, c.cfg.ClientID,
			stan.NatsConn(conn),
			stan.Pings(10, 5),
			stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
				log.Error().Msgf("Connection lost, reason: %v", reason)
			}),
		)

		if err != nil {
			log.Error().Msgf("fail to connect to stan: %s, clusterID: %v, clientID: %v", err.Error(), c.cfg.ClusterID, c.cfg.ClientID)
		}

		err = c.RegisterChannel(context.Background(), c.Channels)
		if err != nil {
			log.Error().Msgf("fail to connect to stan: %s, clusterID: %v, clientID: %v", err.Error(), c.cfg.ClusterID, c.cfg.ClientID)
		}

		return nil
	}, bo)
	if err != nil {
		log.Error().Msgf("fail to reconnect to stan & register channel, err: %s", err.Error())
	}
}

// MsgData ...
type MsgData struct {
	TraceID string          `json:"trace_id"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// Channel ...
type Channel struct {
	ChannelName string
	GroupName   string
	Handler     func(ctx context.Context, data json.RawMessage) error
	Options     []stan.SubscriptionOption
}

// RegisterChannel register channel handler
func RegisterChannel(conn stan.Conn, cfg config.StanConfig, channels []Channel) error {
	for i := range channels {
		logger := log.With().Str("endpoint", channels[i].ChannelName).Logger()
		logger.Info().Msgf("Register channel: %s", channels[i].ChannelName)
		ctx := logger.WithContext(context.Background())

		if cfg.DurableName != "" {
			channels[i].Options = append(channels[i].Options, stan.DurableName(cfg.DurableName))
		}

		_, err := conn.QueueSubscribe(channels[i].ChannelName, channels[i].GroupName, func(msg *stan.Msg) {
			defer recoverLog()

			var msgData MsgData
			err := json.Unmarshal(msg.Data, &msgData)
			if err != nil {
				logger.Error().Msgf("Fail to unmarshal to internal msgData: %s", msgData)
				return
			}

			logger = logger.With().
				Str("trace_id", msgData.TraceID).
				Logger()
			_ctx := logger.WithContext(ctx)
			_ctx = context.WithValue(_ctx, ctxutil.CtxKeyTraceID, msgData.TraceID)

			err = channels[i].Handler(_ctx, msgData.Data)
			if err != nil {
				logger.Error().Msgf("channel: %s, error: %+v", channels[i].ChannelName, err)
			}

		}, channels[i].Options...)
		if err != nil {
			return err
		}
	}
	return nil
}

func recoverLog() {
	if r := recover(); r != nil {
		var msg string
		for i := 2; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			msg += fmt.Sprintf("%s:%d\n", file, line)
		}
		log.Error().Msgf("%s\n↧↧↧↧↧↧ PANIC ↧↧↧↧↧↧\n%s↥↥↥↥↥↥ PANIC ↥↥↥↥↥↥", r, msg)
	}
}
