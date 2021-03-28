package trader

import (
	"context"
	"ptcg_trader/pkg/model"
	"ptcg_trader/pkg/repository"
	"ptcg_trader/pkg/service"

	"go.uber.org/fx"
)

// ServiceParams define params for create service
type ServiceParams struct {
	fx.In

	Repo repository.Repositorier
}

type svc struct {
	repo repository.Repositorier
}

// NewService support DI tool to create a new service instance
func NewService(params ServiceParams) service.TraderServicer {
	return &svc{
		repo: params.Repo,
	}
}

// Get target item by item ID
func (svc *svc) GetItem(ctx context.Context, query model.ItemQuery) (model.Item, error) {
	item, err := svc.repo.GetItem(ctx, query)
	if err != nil {
		return item, err
	}

	return item, nil
}

// List items by query condition
func (svc *svc) ListItems(ctx context.Context, query model.ItemQuery) ([]model.Item, int64, error) {
	items, err := svc.repo.ListItems(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	total, err := svc.repo.CountItems(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// Get target order by order ID
func (svc *svc) GetOrder(ctx context.Context, query model.OrderQuery) (model.Order, error) {
	order, err := svc.repo.GetOrder(ctx, query)
	if err != nil {
		return order, err
	}

	return order, nil
}

// List orders by query condition
func (svc *svc) ListOrders(ctx context.Context, query model.OrderQuery) ([]model.Order, int64, error) {
	orders, err := svc.repo.ListOrders(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	total, err := svc.repo.CountOrders(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Create order
func (svc *svc) CreateOrder(ctx context.Context, order *model.Order) error {
	err := svc.repo.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	// TODO: 嘗試搓合

	return nil
}

// Update orders by query condition
func (svc *svc) UpdateOrders(ctx context.Context, query model.OrderQuery, updates model.OrderUpdates) error {
	err := svc.repo.UpdateOrders(ctx, query, updates)
	if err != nil {
		return err
	}
	return nil
}

// Get target transaction by transaction ID
func (svc *svc) GetTransaction(ctx context.Context, query model.TransactionQuery) (model.Transaction, error) {
	transaction, err := svc.repo.GetTransaction(ctx, query)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// List transactions by query condition
func (svc *svc) ListTransactions(ctx context.Context, query model.TransactionQuery) ([]model.Transaction, int64, error) {
	transactions, err := svc.repo.ListTransactions(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	total, err := svc.repo.CountTransactions(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// Create transaction
func (svc *svc) CreateTransaction(ctx context.Context, tx *model.Transaction) error {
	err := svc.repo.CreateTransaction(ctx, tx)
	if err != nil {
		return err
	}
	return nil
}
