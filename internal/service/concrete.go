package service

import (
	"github.com/SwanHtetAungPhyo/wolftagon/internal/repo"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"sync"
)

type UserService struct {
	log     *logrus.Logger
	repo    *repo.UserRepo
	redis   *redis.Client
	emailWG sync.WaitGroup
}

func NewUserService(
	log *logrus.Logger,
	repo *repo.UserRepo,
	redis *redis.Client,

) *UserService {
	return &UserService{
		log:     log,
		repo:    repo,
		redis:   redis,
		emailWG: sync.WaitGroup{},
	}
}
