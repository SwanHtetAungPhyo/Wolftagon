package repo

import "github.com/SwanHtetAungPhyo/wolftagon/internal/model"

type UserRepoBehaviour interface {
	Create(user *model.User, roleName string) error
	GetByEmail(email string) (*model.User, error)
	MarkAsVerified(email string) error
	UpdatePassword(email string, newPassword string) error
	Delete(email string) error
}
