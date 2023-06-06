package handler

import "github.com/gofiber/fiber/v2"

func RootHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("Hello World!")
	}
}
