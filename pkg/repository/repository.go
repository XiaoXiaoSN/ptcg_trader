package repository

import (
	"context"
	"ptcg_trader/pkg/model"
)

// Repositorier ...
type Repositorier interface {
	TraderRepositorier
	MatcherRepositorier

	// Begin begins a transaction
	Begin(ctx context.Context) Repositorier
	// Commit commit a transaction
	Commit() error
	// Rollback rollback a transaction
	Rollback() error
	// Transaction handle a transaction in a callback function, from begin to commmit
	Transaction(ctx context.Context, f func(ctx context.Context, txRepo Repositorier) error) error
	// Close closes the database and prevents new queries from starting.
	Close() error
}

// TraderRepositorier define PTCG trader service interface
type TraderRepositorier interface {
	// Get target item by item ID
	GetItem(ctx context.Context, query model.ItemQuery) (model.Item, error)
	// counting total count of items
	CountItems(ctx context.Context, query model.ItemQuery) (int64, error)
	// List items by query condition
	ListItems(ctx context.Context, query model.ItemQuery) ([]model.Item, error)

	// Get target order by Order ID
	GetOrder(ctx context.Context, query model.OrderQuery) (model.Order, error)
	// counting total count of orders
	CountOrders(ctx context.Context, query model.OrderQuery) (int64, error)
	// List orders by query condition
	ListOrders(ctx context.Context, query model.OrderQuery) ([]model.Order, error)
	// Create order
	CreateOrder(ctx context.Context, order *model.Order) error
	// Update orders by query condition
	UpdateOrders(ctx context.Context, query model.OrderQuery, updates model.OrderUpdates) error

	// Get target transaction by transaction ID
	GetTransaction(ctx context.Context, query model.TransactionQuery) (model.Transaction, error)
	// counting total count of transactions
	CountTransactions(ctx context.Context, query model.TransactionQuery) (int64, error)
	// List transactions by query condition
	ListTransactions(ctx context.Context, query model.TransactionQuery) ([]model.Transaction, error)
	// Create transaction
	CreateTransaction(ctx context.Context, tx *model.Transaction) error
}

// MatcherRepositorier define PTCG matcher service interface
type MatcherRepositorier interface {
	// MatchOrders check that are there two orders can be matched
	MatchOrders(ctx context.Context, order *model.Order) (model.Order, error)
}
