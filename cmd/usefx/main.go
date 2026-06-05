package main

import (
	"context"
	"time"

	"github.com/PCloud63514/goat/db/mysql"
	"github.com/PCloud63514/goat/environment"
	"github.com/PCloud63514/goat/logger"
	"github.com/PCloud63514/goat/profile"
	"go.uber.org/fx"
)

//////////////////
// Goat Example //
//////////////////

var (
	startUpDateTime = time.Now()
	configurations  = []interface{}{
		&logger.Config{},
		&mysql.Config{},
	}

	goatModule = fx.Module("goat",
		fx.Provide(
			mysql.NewDB,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, db mysql.DB) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						return db.Close()
					},
				})
			},
		),
	)
)

func main() {
	environment.Reset(environment.Option{
		ResPath:        "res",
		Profiles:       profile.Get(),
		Configurations: configurations,
	})

	fx.New(
		fx.Supply(environment.StandardEnvironment()),
		fx.Supply(environment.GetConfigurations()...),
		goatModule,
	).Run()
}
