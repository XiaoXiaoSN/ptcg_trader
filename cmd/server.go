package cmd

import (
	"ptcg_trader/internal/config"
	"ptcg_trader/internal/database"
	"ptcg_trader/internal/echo"
	"ptcg_trader/internal/zlog"
	"ptcg_trader/pkg/delivery/restful"
	"ptcg_trader/pkg/repository/gormrepo"
	"ptcg_trader/pkg/service/trader"

	cobra "github.com/spf13/cobra"
	fx "go.uber.org/fx"
)

// ServerCmd 是此程式的Service入口點
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
			echo.StartEcho,
			gormrepo.NewRepository,
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
