package server_provider

import (
	"github.com/SwanHtetAungPhyo/wolftagon/cmd/server"
	"go.uber.org/fx"
)

var ServerModule = fx.Module("server_provider",
	fx.Provide(server.NewAppState))
