package handler

import "github.com/gofiber/fiber/v2"

type UserHandlerBehaviour interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}
