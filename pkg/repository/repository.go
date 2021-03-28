package repository

import (
	"context"
	"ptcg_trader/pkg/model"
)

// Repositorier ...
type Repositorier interface {
	TraderRepositorier

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
}
