package trader_test

import (
	"context"
	"testing"

	"ptcg_trader/internal/errors"
	"ptcg_trader/pkg/model"
	"ptcg_trader/pkg/service"
	"ptcg_trader/pkg/service/trader"
	"ptcg_trader/test/mocks"

	"github.com/mikunalpha/paws"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

// Test_TraderService define testing case for trader service
func Test_TraderService(t *testing.T) {
	suite.Run(t, &svcTestSuite{})
}

type svcTestSuite struct {
	suite.Suite

	ctx         context.Context
	svc         service.TraderServicer
	mockMatcher *mocks.MockMatcher
	mockRepo    *mocks.MockRepository
	mockRedis   *mocks.MockRedis
}

func (s *svcTestSuite) SetupSuite() {
	s.ctx = context.Background()

	_ = fx.New(
		fx.Provide(
			mocks.ImplMockRepository,
			mocks.ImplMockRedis,
			mocks.ImplMockMatcher,
			trader.NewService,
		),
		fx.Populate(
			&s.mockRepo,
			&s.mockRedis,
			&s.svc,
		),
	)
}

func (s *svcTestSuite) SetupTest() {
}

func (s *svcTestSuite) TearDownTest() {
}

func (s *svcTestSuite) TearDownSuite() {
}

func (s *svcTestSuite) Test_service_GetItem() {
	type args struct {
		ctx   context.Context
		query model.ItemQuery
	}
	tests := []struct {
		name      string
		args      args
		want      model.Item
		wantErr   bool
		prepareFn func(want model.Item, args args)
	}{
		{
			name: "getItem 1",
			args: args{
				ctx: s.ctx,
				query: model.ItemQuery{
					ID: paws.Int64(1),
				}},
			want: model.Item{
				ID:   1,
				Name: "Pikachu",
			},
			wantErr: false,

			prepareFn: func(want model.Item, args args) {
				_ = s.mockRepo.On("GetItem", args.ctx, args.query).
					Return(want, nil).
					Once()
			},
		},
		{
			name: "not found",
			args: args{
				ctx: s.ctx,
				query: model.ItemQuery{
					ID: paws.Int64(5),
				}},
			want:    model.Item{},
			wantErr: true,

			prepareFn: func(want model.Item, args args) {
				_ = s.mockRepo.On("GetItem", args.ctx, args.query).
					Return(model.Item{}, errors.ErrResourceNotFound).
					Once()
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			var err error
			if tt.prepareFn != nil {
				tt.prepareFn(tt.want, tt.args)
			}

			// excute function
			resource, err := s.svc.GetItem(tt.args.ctx, tt.args.query)
			if !tt.wantErr {
				s.Require().NoError(err)
				s.Require().Exactly(resource, tt.want)
			}

			// we make sure that all expectations were met
			isExpectations := s.mockRepo.AssertExpectations(t)
			s.Require().Equal(isExpectations, true)
		})
	}
}
