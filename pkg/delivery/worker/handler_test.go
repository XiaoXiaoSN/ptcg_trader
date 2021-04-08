package worker

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"ptcg_trader/internal/config"
	"ptcg_trader/internal/database"
	"ptcg_trader/internal/pubsub/stan"
	"ptcg_trader/pkg/repository/gormrepo"
	"ptcg_trader/pkg/service"
	"ptcg_trader/pkg/service/matcher"

	"go.uber.org/fx"
)

func TestHandler_CreateOrderEndpoint(t *testing.T) {
	if flag := os.Getenv("E2E_TEST"); flag == "" {
		t.Skip()
		return
	}

	var matcherSvc service.Matcher

	_ = fx.New(
		fx.Provide(
			func() config.Configuration {
				return config.Configuration{
					Database: config.DatabaseConfig{
						Debug:      false,
						Host:       "localhost",
						Port:       5432,
						Username:   "local",
						Password:   "local",
						DBName:     "ptcg",
						Type:       "postgres",
						SearchPath: "trader",
						SSLEnable:  false,
					},
				}
			},
			database.InitDatabases,
			gormrepo.NewRepository,
			matcher.NewMatch,
		),
		fx.Populate(
			&matcherSvc,
		),
	)

	type fields struct {
		stan       *stan.Client
		matcherSvc service.Matcher
	}
	type args struct {
		ctx  context.Context
		data json.RawMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{matcherSvc: matcherSvc},
			args: args{
				ctx: context.Background(),
				data: json.RawMessage(`{
					"id": 1,
					"item_id": 1,
					"order_type": 1,
					"price": "3.8"
				}`),
			},
		},
		{
			fields: fields{matcherSvc: matcherSvc},
			args: args{
				ctx: context.Background(),
				data: json.RawMessage(`{
					"id": 2,
					"item_id": 1,
					"order_type": 1,
					"price": "3.4"
				}`),
			},
		},
		{
			fields: fields{matcherSvc: matcherSvc},
			args: args{
				ctx: context.Background(),
				data: json.RawMessage(`{
					"id": 3,
					"item_id": 1,
					"order_type": 1,
					"price": "3.4"
				}`),
			},
		},
		{
			fields: fields{matcherSvc: matcherSvc},
			args: args{
				ctx: context.Background(),
				data: json.RawMessage(`{
					"id": 4,
					"item_id": 1,
					"order_type": 2,
					"price": "3.5"
				}`),
			},
		},
		{
			fields: fields{matcherSvc: matcherSvc},
			args: args{
				ctx: context.Background(),
				data: json.RawMessage(`{
					"id": 5,
					"item_id": 1,
					"order_type": 2,
					"price": "5"
				}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				stan:       tt.fields.stan,
				matcherSvc: tt.fields.matcherSvc,
			}
			if err := h.CreateOrderEndpoint(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Handler.CreateOrderEndpoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
