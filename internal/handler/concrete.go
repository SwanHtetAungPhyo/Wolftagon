package handler

import (
	"github.com/SwanHtetAungPhyo/wolftagon/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	log   *logrus.Logger
	srv   *service.UserService
	redis *redis.Client
}

func NewUserHandler(
	log *logrus.Logger,
	srv *service.UserService,
	redis *redis.Client,
) *UserHandler {
	return &UserHandler{
		log:   log,
		srv:   srv,
		redis: redis,
	}
}
