package providers

import (
	"github.com/SwanHtetAungPhyo/wolftagon/pkg/database"
	"github.com/SwanHtetAungPhyo/wolftagon/pkg/fiber_app"
	"github.com/SwanHtetAungPhyo/wolftagon/pkg/logs"
	"github.com/SwanHtetAungPhyo/wolftagon/pkg/redis_client"
	"go.uber.org/fx"
)

var PkgModule = fx.Module("pkg", fx.Provide(
	logs.NewLogger,
	redis_client.NewRedisClient,
	database.DBConnect,
	fiber_app.NewFiberApp,
))
