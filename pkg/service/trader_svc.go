package service

import (
	"context"
	"ptcg_trader/pkg/model"
)

// TraderServicer define PTCG trader service interface
type TraderServicer interface {
	// Get target item by item ID
	GetItem(ctx context.Context, query model.ItemQuery) (model.Item, error)
	// List items by query condition
	ListItems(ctx context.Context, query model.ItemQuery) ([]model.Item, int64, error)
}
