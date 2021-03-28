package gormrepo_test

import (
	"context"
	"ptcg_trader/pkg/model"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mikunalpha/paws"
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
				sqlStmt := `SELECT * FROM "items" WHERE id = $1`
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
				sqlStmt := `SELECT * FROM "items" WHERE id = $1`
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
