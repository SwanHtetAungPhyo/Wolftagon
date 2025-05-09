package repo

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepo struct {
	log *logrus.Logger
	db  *gorm.DB
	ctx context.Context
}

func NewUserRepo(log *logrus.Logger, db *gorm.DB) *UserRepo {
	return &UserRepo{log: log, db: db,
		ctx: context.Background()}
}
