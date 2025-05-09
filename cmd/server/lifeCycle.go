package server

import (
	"context"

	"go.uber.org/fx"
)

func RegisterAppLifeCycle(lc fx.Lifecycle, app *AppState) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return app.Start()
		},
		OnStop: func(ctx context.Context) error {
			return app.Stop()
		},
	})
}
