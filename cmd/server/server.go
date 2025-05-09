package server

import (
	"github.com/SwanHtetAungPhyo/wolftagon/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AppState struct {
	log     *logrus.Logger
	db      *gorm.DB
	app     *fiber.App
	handler *handler.UserHandler
	redis   *redis.Client
}

func NewAppState(log *logrus.Logger,
	db *gorm.DB,
	app *fiber.App,
	handler *handler.UserHandler,
	redis *redis.Client,
) *AppState {
	return &AppState{
		log:     log,
		db:      db,
		app:     app,
		handler: handler,
		redis:   redis,
	}
}
