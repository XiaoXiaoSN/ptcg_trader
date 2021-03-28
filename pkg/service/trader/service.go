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
