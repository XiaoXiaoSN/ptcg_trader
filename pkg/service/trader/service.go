package trader

import (
	"ptcg_trader/pkg/repository"
	"ptcg_trader/pkg/service"
)

type svc struct {
	repo repository.Repositorier
}

// NewService support DI tool to create a new service instance
func NewService(repo repository.Repositorier) service.TraderServicer {
	return &svc{repo: repo}
}
