package worker

import (
	"context"
	"encoding/json"
	"ptcg_trader/pkg/model"

	"github.com/rs/zerolog/log"
)

// CreateOrderEndpoint ...
// @Topic  ptcg.order.create
func (h Handler) CreateOrderEndpoint(ctx context.Context, data json.RawMessage) error {
	// revice new order
	order := &model.Order{}
	err := json.Unmarshal(data, order)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("handler: Unmarshal message failed, err: %+v", err)
		return nil
	}

	err = h.matcherSvc.AsyncMatchOrders(ctx, order)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("handler: Async MatchOrders failed, err: %+v", err)
		return nil
	}

	return nil
}
