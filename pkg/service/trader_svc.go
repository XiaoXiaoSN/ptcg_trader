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

	// Get target order by order ID
	GetOrder(ctx context.Context, query model.OrderQuery) (model.Order, error)
	// List orders by query condition
	ListOrders(ctx context.Context, query model.OrderQuery) ([]model.Order, int64, error)
	// Create order
	CreateOrder(ctx context.Context, order *model.Order) error
	// Update orders by query condition
	UpdateOrders(ctx context.Context, query model.OrderQuery, updates model.OrderUpdates) error

	// Get target transaction by transaction ID
	GetTransaction(ctx context.Context, query model.TransactionQuery) (model.Transaction, error)
	// List transactions by query condition
	ListTransactions(ctx context.Context, query model.TransactionQuery) ([]model.Transaction, int64, error)
	// Create transaction
	CreateTransaction(ctx context.Context, tx *model.Transaction) error
}
