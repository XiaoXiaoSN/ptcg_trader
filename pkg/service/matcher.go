package service

import (
	"context"
	"ptcg_trader/pkg/model"
)

// Matcher define PTCG matcher service interface
type Matcher interface {
	// MatchOrder match orders
	MatchOrder(ctx context.Context, order *model.Order) (*model.Order, error)
}
