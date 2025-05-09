package migration

import (
	"github.com/SwanHtetAungPhyo/wolftagon/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, log *logrus.Logger) error {
	if err := db.Exec("SET CONSTRAINTS ALL DEFERRED").Error; err != nil {
		log.Warn("Could not defer constraints")
	}

	tables := []interface{}{
		&model.Role{},
		&model.User{},
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			log.WithError(err).Errorf("Failed to migrate table %T", table)
			return err
		}
	}

	roles := []model.Role{
		{RoleName: "admin"},
		{RoleName: "user"},
		{RoleName: "guest"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, "role_name = ?", role.RoleName).Error; err != nil {
			log.WithError(err).Errorf("Failed to create role %s", role.RoleName)
			return err
		}
	}

	if err := db.Exec("SET CONSTRAINTS ALL IMMEDIATE").Error; err != nil {
		log.Warn("Could not re-enable constraints")
	}

	return nil
}

func DropAllTables(db *gorm.DB, log *logrus.Logger) error {
	if err := db.Exec("SET CONSTRAINTS ALL DEFERRED").Error; err != nil {
		log.Warn("Could not defer constraints")
	}

	tables := []interface{}{
		&model.User{},
		&model.Role{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			log.WithError(err).Errorf("Failed to drop table %T", table)
			return err
		}
	}

	if err := db.Exec("SET CONSTRAINTS ALL IMMEDIATE").Error; err != nil {
		log.Warn("Could not re-enable constraints")
	}

	log.Info("Dropped all database tables")
	return nil
}
