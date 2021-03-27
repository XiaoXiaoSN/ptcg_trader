package repository

import (
	"context"
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
	//
}
