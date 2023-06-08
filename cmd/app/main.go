package main

import (
	"chat_backend/internal/delivery/router"
	"chat_backend/pkg/utils"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gookit/validate"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error reading environment: %v", err)
	}

	app := fiber.New(fiber.Config{
		StrictRouting: true,
		CaseSensitive: true,
		BodyLimit:     10 * 1024 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CLIENT_URL"),
		AllowCredentials: true,
	}))
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(recover.New())

	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})

	db, queries := utils.Database()
	defer db.Close()

	cld, _ := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	router.AppRouter(app, queries, cld)

	log.Fatalln(app.Listen(":6060"))
}
