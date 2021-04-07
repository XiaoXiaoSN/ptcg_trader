package kafka

import (
	"context"
	"encoding/json"

	"ptcg_trader/internal/ctxutil"
	"ptcg_trader/internal/errors"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

// Producer ...
type Producer interface {
	Pub(ctx context.Context, msg PubMsg) error
}

// PubMsg publish msg info
type PubMsg struct {
	// The partitioning key for this message.
	Key   string
	Topic string
	Data  interface{}
}

// SyncProducer ...
type SyncProducer struct {
	producer sarama.SyncProducer
}

// NewSyncProducer ...
func NewSyncProducer(cfg *Config) (*SyncProducer, error) {
	producer, err := newSyncProducer(cfg)
	if err != nil {
		return nil, err
	}
	return &SyncProducer{
		producer: producer,
	}, nil
}

// NewSyncProducerWithFx ...
func NewSyncProducerWithFx(lc fx.Lifecycle, cfg *Config) (Producer, error) {
	sp, err := NewSyncProducer(cfg)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			err := sp.producer.Close()
			if err != nil {
				return err
			}
			return nil
		},
	})

	return sp, nil
}

// Pub ...
func (p *SyncProducer) Pub(ctx context.Context, msg PubMsg) error {
	var msgData MsgData
	var err error
	var b []byte

	switch msg.Data.(type) {
	case []byte:
		b = msg.Data.([]byte)
	default:
		b, err = json.Marshal(msg.Data)
		if err != nil {
			return errors.Wrapf(errors.ErrBadRequest, "fail to marshal pub data, err: %s", err.Error())
		}
	}

	msgData.Data = b
	msgData.TraceID = ctxutil.TraceIDFromCtx(ctx)

	b, err = json.Marshal(msgData)
	if err != nil {
		return errors.Wrapf(errors.ErrBadRequest, "fail to marshal internal MsgData, err: %s", err.Error())
	}

	pMsg := sarama.ProducerMessage{
		Topic: msg.Topic,
		Value: sarama.ByteEncoder(b),
	}
	if len(msg.Key) > 0 {
		pMsg.Key = sarama.StringEncoder(msg.Key)
	}
	partition, offset, err := p.producer.SendMessage(&pMsg)
	if err != nil {
		return errors.Wrapf(errors.ErrInternalError, "fail to send message, err: %s", err.Error())
	}
	log.Ctx(ctx).Debug().Msgf("data is stored with unique identifier important/%d/%d", partition, offset)
	return nil
}

func newSyncProducer(cfg *Config) (sarama.SyncProducer, error) {
	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	config.Producer.CompressionLevel = int(sarama.CompressionZSTD)
	config.Version = sarama.V2_6_0_0

	if cfg.ClientID != "" {
		config.ClientID = cfg.ClientID
	}

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(cfg.BrokerList, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
