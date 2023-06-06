package main

import (
	"chat_backend/internal/delivery/router"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{
		StrictRouting: true,
		CaseSensitive: true,
	})

	router.AppRouter(app)

	log.Fatalln(app.Listen(":6060"))
}
