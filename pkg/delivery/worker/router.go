package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"ptcg_trader/internal/pubsub/stan"
	"ptcg_trader/pkg/model"
	"ptcg_trader/pkg/service"

	pkgStan "github.com/nats-io/stan.go"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func init() {
	// TODO: 本來想做在 stateful-set 的 postStart command
	// https://stackoverflow.com/q/50750672/6695274
	if idx := os.Getenv("INDEX"); idx == "" {
		hostName := os.Getenv("HOSTNAME") // ex: matcher-ptcg-matcher-0
		sp := strings.Split(hostName, "-")
		if len(sp) > 1 {
			err := os.Setenv("INDEX", sp[len(sp)-1])
			fmt.Println("walala:", err)
		}
	}
}

// RegisterChannel register topic
func RegisterChannel(handler *Handler) error {
	ctx := context.Background()

	var indexKeysStr string = "1,2,3,4"
	if idx := os.Getenv("INDEX"); idx != "" {
		indexKeysStr = idx
	}
	indexKeys := strings.Split(indexKeysStr, ",")

	group := newChannelGroup("ptcg")
	{
		var channels = make([]stan.Channel, 0)

		for _, indexKey := range indexKeys {
			// handle model.TopicCreateOrder with index system
			createOrderChannel := fmt.Sprintf("%s.%s", model.TopicCreateOrder, indexKey)
			channels = append(channels, group.register(createOrderChannel, handler.CreateOrderEndpoint))
		}

		err := handler.stan.RegisterChannel(ctx, channels)
		if err != nil {
			log.Error().Msgf("register channels failed: %v", err)
			return err
		}
	}

	return nil
}

func newChannelGroup(group string) chanGroup {
	return chanGroup{GroupName: group}
}

type chanGroup struct {
	GroupName string
}

func (g chanGroup) register(channel string, Handler func(ctx context.Context, data json.RawMessage) error, Options ...pkgStan.SubscriptionOption) stan.Channel {
	return stan.Channel{
		GroupName:   g.GroupName,
		ChannelName: channel,
		Handler:     Handler,
	}
}

// Handler the worker handler
type Handler struct {
	stan *stan.Client

	matcherSvc service.Matcher
}

// Params for worker to new constructor
type Params struct {
	fx.In

	StanClient *stan.Client
	MatcherSvc service.Matcher
}

// NewHandler new a handler for stan worker
func NewHandler(params Params) *Handler {
	return &Handler{
		stan:       params.StanClient,
		matcherSvc: params.MatcherSvc,
	}
}
