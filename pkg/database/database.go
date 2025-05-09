package database

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"sync"
	"time"
)

func DBConnect(log *logrus.Logger) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	var once sync.Once
	dsn := os.Getenv("DATABASE_URL_DNS")

	// Ensure connection happens only once
	once.Do(func() {
		for i := 0; i < 10; i++ {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
			})

			if err == nil {
				// Successfully connected, apply connection pool settings
				rawDB, err := db.DB()
				if err != nil {
					log.Error("Failed to get the raw db connection", err)
					return
				}

				// Connection pool settings
				rawDB.SetMaxIdleConns(10)
				rawDB.SetMaxOpenConns(100)
				rawDB.SetConnMaxLifetime(time.Hour)
				break
			}

			log.Errorf("Failed to connect to the database, attempt %d/%d: %v", i+1, 10, err)

			time.Sleep(2 * time.Second)
		}
	})

	return db, err
}
