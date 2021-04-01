package matcher

import (
	"context"
	"ptcg_trader/pkg/model"
	"ptcg_trader/pkg/repository"
	"ptcg_trader/pkg/service"

	"go.uber.org/fx"
)

// MatchParams define params for create service
type MatchParams struct {
	fx.In

	Repo repository.Repositorier
}

type svc struct {
	repo repository.Repositorier
}

// NewMatch support DI tool to create a new service instance
func NewMatch(params MatchParams) service.Matcher {
	return &svc{
		repo: params.Repo,
	}
}

func (svc *svc) WithRepo(repo repository.Repositorier) service.Matcher {
	if repo == nil {
		repo = svc.repo
	}

	return NewMatch(MatchParams{
		Repo: repo,
	})
}

func (svc *svc) MatchOrders(ctx context.Context, order *model.Order) (*model.Order, error) {
	matchedOrder, err := svc.repo.MatchOrders(ctx, order)
	if err != nil {
		return nil, err
	}

	return &matchedOrder, nil
}
