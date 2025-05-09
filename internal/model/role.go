package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	RoleID   uuid.UUID `gorm:"type:uuid;primaryKey;column:role_id"`
	RoleName string    `gorm:"type:varchar(100);not null;column:role_name"`
}

func (Role) TableName() string {
	return "role_table"
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.RoleID == uuid.Nil {
		r.RoleID = uuid.New()
	}
	return
}
