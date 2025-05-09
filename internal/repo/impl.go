package repo

import (
	"errors"
	"github.com/SwanHtetAungPhyo/wolftagon/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ UserRepoBehaviour = (*UserRepo)(nil)

func (u UserRepo) Create(user *model.User, roleName string) error {
	existingUser, err := u.GetByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		u.log.WithError(err).Error("database error checking user existence")
		return errors.New("database error")
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}

	err = u.db.Transaction(func(tx *gorm.DB) error {
		var role model.Role
		if err := tx.WithContext(u.ctx).
			Where("role_name = ?", roleName).
			FirstOrCreate(&role, model.Role{RoleName: roleName}).Error; err != nil {
			u.log.WithError(err).Error("failed to get/create role")
			return err
		}

		user.RoleID = role.RoleID

		if err := tx.WithContext(u.ctx).Create(user).Error; err != nil {
			u.log.WithError(err).Error("failed to create user")
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	u.log.WithFields(logrus.Fields{
		"user_id": user.UserID,
		"email":   user.Email,
	}).Info("user created successfully")
	return nil
}

func (u UserRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := u.db.WithContext(u.ctx).Preload("Role").First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		u.log.WithError(err).Error("failed to get user by email")
		return nil, err
	}
	return &user, nil
}

func (u UserRepo) UpdatePassword(email string, newPassword string) error {
	result := u.db.WithContext(u.ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Update("password", newPassword)

	if result.Error != nil {
		u.log.WithError(result.Error).Error("failed to update password")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (u UserRepo) Delete(email string) error {
	result := u.db.WithContext(u.ctx).
		Where("email = ?", email).
		Delete(&model.User{})

	if result.Error != nil {
		u.log.WithError(result.Error).Error("failed to delete user")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	u.log.WithField("email", email).Info("user deleted")
	return nil
}

func (u UserRepo) MarkAsVerified(email string) error {
	result := u.db.WithContext(u.ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Update("verified", true)

	if result.Error != nil {
		u.log.WithError(result.Error).Error("failed to mark user as verified")
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (u UserRepo) GetRoleByName(roleName string) (*model.Role, error) {
	var role model.Role
	err := u.db.WithContext(u.ctx).First(&role, "role_name = ?", roleName).Error
	if err != nil {
		u.log.WithError(err).Error("failed to get role by name")
		return nil, err
	}
	return &role, nil
}
