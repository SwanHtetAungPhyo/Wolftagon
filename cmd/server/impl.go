package server

import (
	"context"
	"github.com/SwanHtetAungPhyo/wolftagon/cmd/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"os"
	"time"
)

var _ AppBehaviour = (*AppState)(nil)

func (a *AppState) routeSetup() {
	public := a.app.Group("/auth")
	public.Post("/register", a.handler.Register)
	public.Post("/login", a.handler.Login)
	public.Post("/verify", a.handler.Verify)

	a.app.Post("/logout", middleware.JwtMiddleware(a.redis, a.log), a.handler.Logout)
	a.app.Get("/refresh", middleware.JwtMiddleware(a.redis, a.log), a.handler.Refresh)
	admin := a.app.Group("/admin")
	admin.Use(middleware.RoleMiddleware(a.redis, a.log, "admin"))
	admin.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"success": true,
			"message": "Welcome  Wolftagon Addmin",
		})
	})
	user := a.app.Group("/user")
	user.Use(middleware.RoleMiddleware(a.redis, a.log, "user"))
	user.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"success": true,
			"message": "Welcome  Wolftagon User",
		})
	})
}

func (a *AppState) middlewareSetup() {
	a.app.Use(recover.New())

	a.app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("x-forwarded-for", c.IP())
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "rate limit exceeded",
			})
		},
	}))

}
func (a *AppState) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8081"
	}
	pwd, _ := os.Getwd()

	certFile := pwd + "/cmd/server" + os.Getenv("SSL_CERT_FILE")
	keyFile := pwd + "/cmd/server" + os.Getenv("SSL_KEY_FILE")
	a.middlewareSetup()
	a.routeSetup()
	a.log.Info("application started")
	var err error
	go func() {
		if err = a.app.ListenTLS(port, certFile, keyFile); err != nil {
			return
		}
	}()
	return nil
}

func (a *AppState) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return a.app.ShutdownWithContext(ctx)
}
