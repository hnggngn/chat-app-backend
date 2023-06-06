package router

import (
	"chat_backend/internal/delivery/handler"
	"github.com/gofiber/fiber/v2"
)

func AppRouter(app *fiber.App) {
	api := app.Group("/api")

	api.Get("", handler.RootHandler())
}
