package cmd

import (
	"ptcg_trader/internal/config"
	"ptcg_trader/internal/database"
	"ptcg_trader/internal/echo"
	"ptcg_trader/internal/pubsub/stan"
	"ptcg_trader/internal/redis"
	"ptcg_trader/internal/zlog"
	"ptcg_trader/pkg/delivery/restful"
	"ptcg_trader/pkg/repository/gormrepo"
	"ptcg_trader/pkg/service/matcher"
	"ptcg_trader/pkg/service/trader"

	cobra "github.com/spf13/cobra"
	fx "go.uber.org/fx"
)

// ServerCmd is the service enteripoint
var ServerCmd = &cobra.Command{
	RunE: runServerCmd,
	Use:  "server",
}

func runServerCmd(cmd *cobra.Command, args []string) error {
	defer cmdRecover()

	app := fx.New(
		fx.Provide(
			config.New,
			database.InitDatabases,
			redis.NewRedis,
			echo.StartEcho,
			gormrepo.NewRepository,
			stan.NewClientWithFx,
			matcher.NewMatch,
			trader.NewService,
			restful.NewHandler,
		),
		fx.Invoke(
			zlog.InitLog,
			restful.SetRoutes,
		),
	)

	return startGraceApp(cmd.Name(), app)
}
