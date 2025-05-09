package internal_provider

import (
	"github.com/SwanHtetAungPhyo/wolftagon/internal/handler"
	"github.com/SwanHtetAungPhyo/wolftagon/internal/repo"
	"github.com/SwanHtetAungPhyo/wolftagon/internal/service"
	"go.uber.org/fx"
)

var InternalModule = fx.Module("internal_provider",
	fx.Provide(
		repo.NewUserRepo,
		service.NewUserService,
		handler.NewUserHandler,
	))
