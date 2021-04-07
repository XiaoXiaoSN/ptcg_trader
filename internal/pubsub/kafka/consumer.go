package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"ptcg_trader/internal/ctxutil"

	"github.com/Shopify/sarama"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

// TopicName ...
type TopicName string

// HandleFunc ...
type HandleFunc func(ctx context.Context, data json.RawMessage) error

// Consumer ...
type Consumer interface {
	Stop() error
	SetupConsumerHandler(cgtc *ConsumerGroupTopicsHandlerConfig)
}

// ConsumerGroupTopicsHandlerConfig ...
type ConsumerGroupTopicsHandlerConfig struct {
	Handler map[TopicName]HandleFunc
}

// ConsumerGroup represents a Sarama consumer group consumer
type ConsumerGroup struct {
	cfg           *Config
	ready         chan bool
	wg            *sync.WaitGroup
	client        sarama.ConsumerGroup
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
	cgtc          *ConsumerGroupTopicsHandlerConfig
}

// NewConsumerGroup ...
func NewConsumerGroup(cfg *Config) (Consumer, error) {
	client, err := newConsumerGroup(cfg)
	if err != nil {
		return nil, err
	}
	consumerGroup := &ConsumerGroup{
		cfg:   cfg,
		ready: make(chan bool),
		wg:    &sync.WaitGroup{},
	}
	ctx, cancel := context.WithCancel(context.Background())
	consumerGroup.ctx = ctx
	consumerGroup.ctxCancelFunc = cancel
	consumerGroup.client = client

	return consumerGroup, nil
}

// NewConsumerWithFx ...
func NewConsumerWithFx(lc fx.Lifecycle, cfg *Config) (Consumer, error) {
	cg, err := NewConsumerGroup(cfg)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			err := cg.Stop()
			if err != nil {
				return err
			}
			return nil
		},
	})

	return cg, nil
}

func newConsumerGroup(cfg *Config) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	if cfg.OffsetsInitial == 0 {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		config.Consumer.Offsets.Initial = cfg.OffsetsInitial
	}
	if cfg.ClientID != "" {
		config.ClientID = cfg.ClientID
	}

	client, err := sarama.NewConsumerGroup(cfg.BrokerList, cfg.GroupName, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Stop stop consumer
func (consumer *ConsumerGroup) Stop() error {
	consumer.ctxCancelFunc()
	consumer.wg.Wait()
	err := consumer.client.Close()
	if err != nil {
		return err
	}
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *ConsumerGroup) Setup(session sarama.ConsumerGroupSession) error {
	for topic, partitions := range session.Claims() {
		log.Debug().Msgf("consumer: listen on topic(%s), partitions(%+v)", topic, partitions)
	}

	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		handler, ok := consumer.cgtc.Handler[TopicName(message.Topic)]
		if !ok {
			continue
		}

		logger := log.With().Str("endpoint", message.Topic).Logger()

		var msgData MsgData
		err := json.Unmarshal(message.Value, &msgData)
		if err != nil {
			logger.Error().Msgf("Fail to unmarshal to internal msgData: %s", msgData)
			continue
		}
		msgData.ConsumeID = xid.New().String()

		logger = logger.With().
			Str("trace_id", msgData.TraceID).
			Str("consume_id", msgData.ConsumeID).
			Logger()

		if consumer.cfg.Debug {
			logger.
				Debug().
				Str("kafka_key", string(message.Key)).
				Str("kafka_topic", message.Topic).
				Int32("kafka_offset", int32(message.Offset)).
				Int32("kafka_partition", message.Partition).
				Str("kafka_group", consumer.cfg.GroupName).
				Str("kafka_req", string(msgData.Data)).
				Str("kafka_msg_time", message.Timestamp.String()).
				Msg("access log")
		}

		ctx := logger.WithContext(context.Background())
		ctx = context.WithValue(ctx, ctxutil.CtxKeyTraceID, msgData.TraceID)

		err = consumer.execHandler(ctx, handler, msgData.Data, logger)
		if err != nil {
			logger.Error().Msgf("channel: %s, message: %+v , error: %+v", message.Topic, message, err)
			continue
		}

		session.MarkMessage(message, "")
	}

	return nil
}

func (consumer *ConsumerGroup) execHandler(ctx context.Context, handler HandleFunc, data json.RawMessage, logger zerolog.Logger) error {
	defer func() {
		if r := recover(); r != nil {
			var msg string
			for i := 2; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				msg += fmt.Sprintf("%s:%d\n", file, line)
			}
			logger.Error().Msgf("%s\n↧↧↧↧↧↧ PANIC ↧↧↧↧↧↧\n%s↥↥↥↥↥↥ PANIC ↥↥↥↥↥↥", r, msg)
		}
	}()
	return handler(ctx, data)
}

// SetupConsumerHandler ...
func (consumer *ConsumerGroup) SetupConsumerHandler(cgtc *ConsumerGroupTopicsHandlerConfig) {
	consumer.cgtc = cgtc
	topicList := make([]string, len(cgtc.Handler))
	idx := 0
	for k := range cgtc.Handler {
		topicList[idx] = string(k)
		idx++
	}
	consumer.wg.Add(1)
	go func() {
		defer consumer.wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := consumer.client.Consume(consumer.ctx, topicList, consumer); err != nil {
				log.Error().Err(err).Msgf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if consumer.ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Debug().Msg("Sarama consumer up and running!...")
}
