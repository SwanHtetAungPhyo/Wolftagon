package main

import (
	"github.com/SwanHtetAungPhyo/wolftagon/cmd/server"
	"github.com/SwanHtetAungPhyo/wolftagon/cmd/server/migration"
	"github.com/SwanHtetAungPhyo/wolftagon/cmd/server/server_provider"
	internal_provider "github.com/SwanHtetAungPhyo/wolftagon/internal/provider"
	"github.com/SwanHtetAungPhyo/wolftagon/pkg/providers"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	pwd, _ := os.Getwd()
	filPath := filepath.Join(pwd, "/cmd", ".env")
	err := godotenv.Load(filPath)
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

	app := fx.New(
		providers.PkgModule,
		internal_provider.InternalModule,
		server_provider.ServerModule,
		fx.Invoke(
			migration.Migrate,
			server.RegisterAppLifeCycle,
		),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		app.Run()
	}()

	sig := <-sigChan
	log.Infof("Received signal: %v", sig)

	var (
		db  *gorm.DB
		log *logrus.Logger
	)

	if os.Getenv("ENVIRONMENT") == "development" {
		log.Info("Cleaning up database...")
		if err := migration.DropAllTables(db, log); err != nil {
			log.Errorf("Failed to drop tables: %v", err)
		}
	}

	app.Done()
}
