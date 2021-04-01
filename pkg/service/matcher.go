package service

import (
	"context"
	"ptcg_trader/pkg/model"
	"ptcg_trader/pkg/repository"
)

// Matcher define PTCG matcher service interface
type Matcher interface {
	// WithRepo return a new Matcher with the designed repository
	WithRepo(repo repository.Repositorier) Matcher

	// MatchOrders match orders
	MatchOrders(ctx context.Context, order *model.Order) (*model.Order, error)
}
