package gormrepo_test

import (
	"context"
	"ptcg_trader/pkg/model"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mikunalpha/paws"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func (s *repoTestSuite) Test_repository_GetItem() {
	type args struct {
		ctx   context.Context
		query model.ItemQuery
	}
	tests := []struct {
		name      string
		args      args
		want      model.Item
		wantErr   bool
		prepareFn func(want model.Item, query model.ItemQuery)
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

			prepareFn: func(want model.Item, query model.ItemQuery) {
				row := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(want.ID, want.Name)
				sqlStmt := `SELECT * FROM "items" WHERE "items"."id" = $1`
				s.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlStmt)).
					WithArgs(want.ID).
					WillReturnRows(row)
			},
		},
		{
			name: "getItem not found",
			args: args{
				ctx: s.ctx,
				query: model.ItemQuery{
					ID: paws.Int64(5),
				}},
			want:    model.Item{},
			wantErr: true,

			prepareFn: func(want model.Item, query model.ItemQuery) {
				sqlStmt := `SELECT * FROM "items" WHERE "items"."id" = $1`
				s.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlStmt)).
					WithArgs(query.ID).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
	}
	for _, tt := range tests {
		var err error
		if tt.prepareFn != nil {
			tt.prepareFn(tt.want, tt.args.query)
		}

		// excute function
		resource, err := s.Repo.GetItem(tt.args.ctx, tt.args.query)
		if !tt.wantErr {
			s.Require().NoError(err)
			s.Require().Exactly(resource, tt.want)
		}

		// we make sure that all expectations were met
		err = s.SQLMock.ExpectationsWereMet()
		s.Require().NoError(err)
	}
}

func (s *repoTestSuite) Test_repository_CountItems() {
	type args struct {
		ctx   context.Context
		query model.ItemQuery
	}
	tests := []struct {
		name      string
		args      args
		want      int64
		wantErr   bool
		prepareFn func(want int64, query model.ItemQuery)
	}{
		{
			name: "count item",
			args: args{
				ctx:   s.ctx,
				query: model.ItemQuery{
					//
				}},
			want:    10,
			wantErr: false,

			prepareFn: func(want int64, query model.ItemQuery) {
				row := sqlmock.NewRows([]string{"count"}).AddRow(want)
				sqlStmt := `SELECT count(1) FROM "items"`
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
		resource, err := s.Repo.CountItems(tt.args.ctx, tt.args.query)
		if !tt.wantErr {
			s.Require().NoError(err)
			s.Require().Exactly(resource, tt.want)
		}

		// we make sure that all expectations were met
		err = s.SQLMock.ExpectationsWereMet()
		s.Require().NoError(err)
	}
}

func (s *repoTestSuite) Test_repository_ListItems() {
	type args struct {
		ctx   context.Context
		query model.ItemQuery
	}
	tests := []struct {
		name      string
		args      args
		want      []model.Item
		wantErr   bool
		prepareFn func(want []model.Item, query model.ItemQuery)
	}{
		{
			name: "list items",
			args: args{
				ctx:   s.ctx,
				query: model.ItemQuery{
					//
				}},
			want: []model.Item{
				{ID: 1, Name: "pika", ImageURL: "http://", CreatorID: 1},
				{ID: 2, Name: "pika 2", ImageURL: "https://", CreatorID: 2},
			},
			wantErr: false,

			prepareFn: func(want []model.Item, query model.ItemQuery) {
				rows := sqlmock.NewRows([]string{"id", "name", "image_url", "creator_id"})
				for _, r := range want {
					rows.AddRow(r.ID, r.Name, r.ImageURL, r.CreatorID)
				}
				sqlStmt := `SELECT * FROM "items"`
				s.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlStmt)).
					WillReturnRows(rows)
			},
		},
	}
	for _, tt := range tests {
		var err error
		if tt.prepareFn != nil {
			tt.prepareFn(tt.want, tt.args.query)
		}

		// excute function
		resource, err := s.Repo.ListItems(tt.args.ctx, tt.args.query)
		if !tt.wantErr {
			s.Require().NoError(err)
			s.Require().Exactly(resource, tt.want)
		}

		// we make sure that all expectations were met
		err = s.SQLMock.ExpectationsWereMet()
		s.Require().NoError(err)
	}
}

func (s *repoTestSuite) Test_repository_GetOrder() {
	type args struct {
		ctx   context.Context
		query model.OrderQuery
	}
	tests := []struct {
		name      string
		args      args
		want      model.Order
		wantErr   bool
		prepareFn func(want model.Order, query model.OrderQuery)
	}{
		{
			name: "getOrder 1",
			args: args{
				ctx: s.ctx,
				query: model.OrderQuery{
					ID: paws.Int64(1),
				}},
			want: model.Order{
				ID:        1,
				ItemID:    1,
				CreatorID: 2,
				OrderType: model.OrderTypeBuy,
				Price:     decimal.NewFromFloat(8.10),
				Status:    model.OrderStatusCompleted,
			},
			wantErr: false,

			prepareFn: func(want model.Order, query model.OrderQuery) {
				row := sqlmock.
					NewRows([]string{"id", "item_id", "creator_id", "order_type", "price", "status"}).
					AddRow(want.ID, want.ItemID, want.CreatorID, want.OrderType, want.Price, want.Status)
				sqlStmt := `SELECT * FROM "orders" WHERE "orders"."id" = $1`
				s.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlStmt)).
					WithArgs(want.ID).
					WillReturnRows(row)
			},
		},
		{
			name: "getOrder not found",
			args: args{
				ctx: s.ctx,
				query: model.OrderQuery{
					ID: paws.Int64(5),
				}},
			want:    model.Order{},
			wantErr: true,

			prepareFn: func(want model.Order, query model.OrderQuery) {
				sqlStmt := `SELECT * FROM "orders" WHERE "orders"."id" = $1`
				s.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlStmt)).
					WithArgs(query.ID).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
	}
	for _, tt := range tests {
		var err error
		if tt.prepareFn != nil {
			tt.prepareFn(tt.want, tt.args.query)
		}

		// excute function
		resource, err := s.Repo.GetOrder(tt.args.ctx, tt.args.query)
		if !tt.wantErr {
			s.Require().NoError(err)
			s.Require().Exactly(resource, tt.want)
		}

		// we make sure that all expectations were met
		err = s.SQLMock.ExpectationsWereMet()
		s.Require().NoError(err)
	}
}

func (s *repoTestSuite) Test_repository_CountOrders() {
	type args struct {
		ctx   context.Context
		query model.OrderQuery
	}
	tests := []struct {
		name      string
		args      args
		want      int64
		wantErr   bool
		prepareFn func(want int64, query model.OrderQuery)
	}{
		{
			name: "count order",
			args: args{
				ctx:   s.ctx,
				query: model.OrderQuery{
					//
				}},
			want:    10,
			wantErr: false,

			prepareFn: func(want int64, query model.OrderQuery) {
				row := sqlmock.NewRows([]string{"count"}).AddRow(want)
				sqlStmt := `SELECT count(1) FROM "orders"`
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
		resource, err := s.Repo.CountOrders(tt.args.ctx, tt.args.query)
		if !tt.wantErr {
			s.Require().NoError(err)
			s.Require().Exactly(resource, tt.want)
		}

		// we make sure that all expectations were met
		err = s.SQLMock.ExpectationsWereMet()
		s.Require().NoError(err)
	}
}

func (s *repoTestSuite) Test_repository_ListOrders() {
	type args struct {
		ctx   context.Context
		query model.OrderQuery
	}
	tests := []struct {
		name      string
		args      args
		want      []model.Order
		wantErr   bool
		prepareFn func(want []model.Order, query model.OrderQuery)
	}{
		{
			name: "list orders",
			args: args{
				ctx:   s.ctx,
				query: model.OrderQuery{
					//
				}},
			want: []model.Order{
				{
					ID:        1,
					ItemID:    1,
					CreatorID: 2,
					OrderType: model.OrderTypeBuy,
					Price:     decimal.NewFromFloat(8.10),
					Status:    model.OrderStatusCompleted,
				},
				{
					ID:        2,
					ItemID:    2,
					CreatorID: 2,
					OrderType: model.OrderTypeSell,
					Price:     decimal.NewFromFloat(7.10),
					Status:    model.OrderStatusCompleted,
				},
			},
			wantErr: false,

			prepareFn: func(want []model.Order, query model.OrderQuery) {
				rows := sqlmock.NewRows([]string{"id", "item_id", "creator_id", "order_type", "price", "status"})
				for _, r := range want {
					rows.AddRow(r.ID, r.ItemID, r.CreatorID, r.OrderType, r.Price, r.Status)
				}
				sqlStmt := `SELECT * FROM "orders"`
				s.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlStmt)).
					WillReturnRows(rows)
			},
		},
	}
	for _, tt := range tests {
		var err error
		if tt.prepareFn != nil {
			tt.prepareFn(tt.want, tt.args.query)
		}

		// excute function
		resource, err := s.Repo.ListOrders(tt.args.ctx, tt.args.query)
		if !tt.wantErr {
			s.Require().NoError(err)
			s.Require().Exactly(resource, tt.want)
		}

		// we make sure that all expectations were met
		err = s.SQLMock.ExpectationsWereMet()
		s.Require().NoError(err)
	}
}

func (s *repoTestSuite) Test_repository_GetTransaction() {
	type args struct {
		ctx   context.Context
		query model.TransactionQuery
	}
	tests := []struct {
		name      string
		args      args
		want      model.Transaction
		wantErr   bool
		prepareFn func(want model.Transaction, query model.TransactionQuery)
	}{
		{
			name: "getTransaction 1",
			args: args{
				ctx: s.ctx,
				query: model.TransactionQuery{
					ID: paws.Int64(1),
				}},
			want: model.Transaction{
				ID:          1,
				MakeOrderID: 1,
				TakeOrderID: 2,
				FinalPrice:  decimal.NewFromFloat(9.0),
			},
			wantErr: false,

			prepareFn: func(want model.Transaction, query model.TransactionQuery) {
				row := sqlmock.NewRows([]string{"id", "make_order_id", "take_order_id", "final_price"}).
					AddRow(want.ID, want.MakeOrderID, want.TakeOrderID, want.FinalPrice)
				sqlStmt := `SELECT * FROM "transactions" WHERE "transactions"."id" = $1`
				s.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlStmt)).
					WithArgs(want.ID).
					WillReturnRows(row)
			},
		},
		{
			name: "getTransaction not found",
			args: args{
				ctx: s.ctx,
				query: model.TransactionQuery{
					ID: paws.Int64(5),
				}},
			want:    model.Transaction{},
			wantErr: true,

			prepareFn: func(want model.Transaction, query model.TransactionQuery) {
				sqlStmt := `SELECT * FROM "transactions" WHERE "transactions"."id" = $1`
				s.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlStmt)).
					WithArgs(query.ID).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
	}
	for _, tt := range tests {
		var err error
		if tt.prepareFn != nil {
			tt.prepareFn(tt.want, tt.args.query)
		}

		// excute function
		resource, err := s.Repo.GetTransaction(tt.args.ctx, tt.args.query)
		if !tt.wantErr {
			s.Require().NoError(err)
			s.Require().Exactly(resource, tt.want)
		}

		// we make sure that all expectations were met
		err = s.SQLMock.ExpectationsWereMet()
		s.Require().NoError(err)
	}
}

func (s *repoTestSuite) Test_repository_CountTransactions() {
	type args struct {
		ctx   context.Context
		query model.TransactionQuery
	}
	tests := []struct {
		name      string
		args      args
		want      int64
		wantErr   bool
		prepareFn func(want int64, query model.TransactionQuery)
	}{
		{
			name: "count transaction",
			args: args{
				ctx:   s.ctx,
				query: model.TransactionQuery{
					//
				}},
			want:    10,
			wantErr: false,

			prepareFn: func(want int64, query model.TransactionQuery) {
				row := sqlmock.NewRows([]string{"count"}).AddRow(want)
				sqlStmt := `SELECT count(1) FROM "transactions"`
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
		resource, err := s.Repo.CountTransactions(tt.args.ctx, tt.args.query)
		if !tt.wantErr {
			s.Require().NoError(err)
			s.Require().Exactly(resource, tt.want)
		}

		// we make sure that all expectations were met
		err = s.SQLMock.ExpectationsWereMet()
		s.Require().NoError(err)
	}
}

func (s *repoTestSuite) Test_repository_ListTransactions() {
	type args struct {
		ctx   context.Context
		query model.TransactionQuery
	}
	tests := []struct {
		name      string
		args      args
		want      []model.Transaction
		wantErr   bool
		prepareFn func(want []model.Transaction, query model.TransactionQuery)
	}{
		{
			name: "list transactions",
			args: args{
				ctx:   s.ctx,
				query: model.TransactionQuery{
					//
				}},
			want: []model.Transaction{
				{
					ID:          1,
					MakeOrderID: 1,
					TakeOrderID: 2,
					FinalPrice:  decimal.NewFromFloat(9.0),
				},
				{
					ID:          2,
					MakeOrderID: 3,
					TakeOrderID: 4,
					FinalPrice:  decimal.NewFromFloat(8.0),
				},
			},
			wantErr: false,

			prepareFn: func(want []model.Transaction, query model.TransactionQuery) {
				rows := sqlmock.NewRows([]string{"id", "make_order_id", "take_order_id", "final_price"})
				for _, r := range want {
					rows.AddRow(r.ID, r.MakeOrderID, r.TakeOrderID, r.FinalPrice)
				}
				sqlStmt := `SELECT * FROM "transactions"`
				s.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlStmt)).
					WillReturnRows(rows)
			},
		},
	}
	for _, tt := range tests {
		var err error
		if tt.prepareFn != nil {
			tt.prepareFn(tt.want, tt.args.query)
		}

		// excute function
		resource, err := s.Repo.ListTransactions(tt.args.ctx, tt.args.query)
		if !tt.wantErr {
			s.Require().NoError(err)
			s.Require().Exactly(resource, tt.want)
		}

		// we make sure that all expectations were met
		err = s.SQLMock.ExpectationsWereMet()
		s.Require().NoError(err)
	}
}
