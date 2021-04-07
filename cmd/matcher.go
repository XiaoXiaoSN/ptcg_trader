package cmd

import (
	"ptcg_trader/internal/config"
	"ptcg_trader/internal/database"
	"ptcg_trader/internal/pubsub/stan"
	"ptcg_trader/internal/redis"
	"ptcg_trader/internal/zlog"
	"ptcg_trader/pkg/delivery/worker"
	"ptcg_trader/pkg/repository/gormrepo"
	"ptcg_trader/pkg/service/matcher"

	cobra "github.com/spf13/cobra"
	fx "go.uber.org/fx"
)

// ServerCmd is the service enteripoint
var MatcherCmd = &cobra.Command{
	RunE: runMatcherCmd,
	Use:  "matcher",
}

func runMatcherCmd(cmd *cobra.Command, args []string) error {
	defer cmdRecover()

	app := fx.New(
		fx.Provide(
			config.New,
			database.InitDatabases,
			redis.NewRedis,
			gormrepo.NewRepository,
			stan.NewClientWithFx,
			matcher.NewMatch,
			worker.NewHandler,
		),
		fx.Invoke(
			zlog.InitLog,
			worker.RegisterChannel,
		),
	)

	return startGraceApp(cmd.Name(), app)
}
