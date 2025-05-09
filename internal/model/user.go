package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	FirstName string    `gorm:"type:varchar(100);not null;column:first_name"`
	LastName  string    `gorm:"type:varchar(100);not null;column:last_name"`
	Age       int       `gorm:"not null"`
	Password  string    `gorm:"not null"`
	Verified  bool      `gorm:"not null;column:verified"`
	Email     string    `gorm:"type:varchar(150);not null;uniqueIndex"`
	RoleID    uuid.UUID `gorm:"type:uuid;not null;column:role_id"`
	Role      Role      `gorm:"foreignKey:RoleID;references:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (User) TableName() string {
	return "user_table"
}
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}
	return
}
