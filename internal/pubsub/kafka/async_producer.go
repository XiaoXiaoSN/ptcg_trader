package kafka

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"ptcg_trader/internal/ctxutil"
	"ptcg_trader/internal/errors"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

// AsyncProducer ...
type AsyncProducer interface {
	Pub(ctx context.Context, msg PubMsg) error
	SetAsyncProducerErrorCB(cb AsyncProducerErrorHandler)
	SetAsyncProducerSuccessCB(cb AsyncProducerSuccessHandler)
	AsyncClose()
}

// AsyncProducerErrorHandler ...
type AsyncProducerErrorHandler func(ProducerError)

// AsyncProducerSuccessHandler ...
type AsyncProducerSuccessHandler func(ProducerMessage)

// AsyncProducerOptions ...
type AsyncProducerOptions struct {
	AsyncProducerErrorCB   AsyncProducerErrorHandler
	AsyncProducerSuccessCB AsyncProducerSuccessHandler
}

type asyncProducer struct {
	producer sarama.AsyncProducer
	wg       sync.WaitGroup
	opts     AsyncProducerOptions
}

// NewAsyncProducer ...
func NewAsyncProducer(cfg *Config) (AsyncProducer, error) {
	producer, err := newAsyncProducer(cfg)
	if err != nil {
		return nil, err
	}
	ap := &asyncProducer{
		producer: producer,
	}

	ap.wg.Add(1)
	go func() {
		defer ap.wg.Done()
		for msg := range ap.producer.Successes() {
			if ap.opts.AsyncProducerSuccessCB != nil {
				data, err := msg.Value.Encode()
				if err != nil {
					log.Error().Msgf("encode message failed: %v", err)
					continue
				}

				var msgData MsgData
				err = json.Unmarshal(data, &msgData)
				if err != nil {
					log.Error().Msgf("in producer Successes, json unmarshal data to MsgData failed: %v", err)
					continue
				}

				producerMessage := ProducerMessage{
					Topic: msg.Topic,
					Data:  msgData.Data,
				}
				ap.opts.AsyncProducerSuccessCB(producerMessage)
			}
		}
	}()

	ap.wg.Add(1)
	go func() {
		defer ap.wg.Done()
		for pErr := range ap.producer.Errors() {
			if ap.opts.AsyncProducerSuccessCB != nil {
				data, err := pErr.Msg.Value.Encode()
				if err != nil {
					log.Error().Msgf("encode message failed: %v", err)
					continue
				}

				var msgData MsgData
				err = json.Unmarshal(data, &msgData)
				if err != nil {
					log.Error().Msgf("in producer Errors, json unmarshal data to MsgData failed: %v", err)
					continue
				}

				producerError := ProducerError{
					Msg: ProducerMessage{
						Topic: pErr.Msg.Topic,
						Data:  data,
					},
					Err: pErr.Err,
				}
				ap.opts.AsyncProducerErrorCB(producerError)
			}
		}
	}()

	return ap, nil
}

// NewAsyncProducerWithFx ...
func NewAsyncProducerWithFx(lc fx.Lifecycle, cfg *Config) (AsyncProducer, error) {
	ap, err := NewAsyncProducer(cfg)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			ap.AsyncClose()
			return nil
		},
	})

	return ap, nil
}

// SetAsyncProducerErrorCB ...
func (p *asyncProducer) SetAsyncProducerErrorCB(cb AsyncProducerErrorHandler) {
	p.opts.AsyncProducerErrorCB = cb
}

// SetAsyncProducerSuccessCB ...
func (p *asyncProducer) SetAsyncProducerSuccessCB(cb AsyncProducerSuccessHandler) {
	p.opts.AsyncProducerSuccessCB = cb
}

// Pub ...
func (p *asyncProducer) Pub(ctx context.Context, msg PubMsg) error {
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
	p.producer.Input() <- &pMsg
	return nil
}

// AsyncClose ...
func (p *asyncProducer) AsyncClose() {
	p.producer.AsyncClose()
	p.wg.Wait()
}

func newAsyncProducer(cfg *Config) (sarama.AsyncProducer, error) {

	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionZSTD     // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	config.Version = sarama.V2_6_0_0

	if cfg.ClientID != "" {
		config.ClientID = cfg.ClientID
	}

	producer, err := sarama.NewAsyncProducer(cfg.BrokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
