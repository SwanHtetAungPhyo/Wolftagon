package fiber_app

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func NewFiberApp(log *logrus.Logger) *fiber.App {
	idleTimeoutStr := os.Getenv("APP_IDLE_TIMEOUT")
	readTimeoutStr := os.Getenv("APP_READ_TIMEOUT")
	writeTimeoutStr := os.Getenv("APP_WRITE_TIMEOUT")

	idleTimeout, err := time.ParseDuration(idleTimeoutStr)
	if err != nil {
		log.Warnf("Invalid APP_IDLE_TIMEOUT: %v. Using default 60s.", err)
		idleTimeout = 60 * time.Second
	}

	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		log.Warnf("Invalid APP_READ_TIMEOUT: %v. Using default 10s.", err)
		readTimeout = 10 * time.Second
	}

	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		log.Warnf("Invalid APP_WRITE_TIMEOUT: %v. Using default 10s.", err)
		writeTimeout = 10 * time.Second
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		IdleTimeout:           idleTimeout,
		ReadTimeout:           readTimeout,
		WriteTimeout:          writeTimeout,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
		AppName:           "fiber_user_app",
		ReduceMemoryUsage: true,
		EnablePrintRoutes: true,
	})

	return app
}
