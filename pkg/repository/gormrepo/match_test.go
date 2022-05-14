package gormrepo_test

import (
	"context"
	"regexp"

	"ptcg_trader/pkg/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
)

func (s *repoTestSuite) Test_repository_MatchOrders() {
	type args struct {
		ctx   context.Context
		query *model.Order
	}
	tests := []struct {
		name      string
		args      args
		want      model.Order
		wantErr   bool
		prepareFn func(want model.Order, query *model.Order)
	}{
		{
			name: "match orders",
			args: args{
				ctx: s.ctx,
				query: &model.Order{
					ID:        2,
					ItemID:    1,
					OrderType: model.OrderTypeBuy,
					Price:     decimal.NewFromFloat(5.5),
					Status:    model.OrderStatusProgress,
				}},
			want: model.Order{
				ID:        1,
				ItemID:    1,
				OrderType: model.OrderTypeSell,
				Price:     decimal.NewFromFloat(5.0),
				Status:    model.OrderStatusProgress,
			},
			wantErr: false,

			prepareFn: func(want model.Order, query *model.Order) {
				row := sqlmock.NewRows([]string{"id", "item_id", "order_type", "price", "status"}).
					AddRow(want.ID, want.ItemID, want.OrderType, want.Price, want.Status)
				sqlStmt := `SELECT * FROM "orders"`
				s.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlStmt)).
					WillReturnRows(row)
			},
		},
	}
	for _, tt := range tests {
		var err error
		if tt.prepareFn != nil {
			tt.prepareFn(tt.want, tt.args.query)
		}

		// excute function
		resource, err := s.Repo.MatchOrders(tt.args.ctx, tt.args.query)
		if !tt.wantErr {
			s.Require().NoError(err)
			s.Require().Exactly(resource, tt.want)
		}

		// we make sure that all expectations were met
		err = s.SQLMock.ExpectationsWereMet()
		s.Require().NoError(err)
	}
}
