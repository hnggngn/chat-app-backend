package router

import (
	"chat_backend/generated"
	"chat_backend/internal/app/repositories"
	"chat_backend/internal/app/services"
	"chat_backend/internal/delivery/handlers"
	"chat_backend/pkg/utils"
	pasetoware "github.com/gofiber/contrib/paseto"
	"github.com/gofiber/fiber/v2"
)

func AppRouter(app *fiber.App, queries *generated.Queries) {
	authRepo := repositories.NewRepo(queries)
	authService := services.NewService(authRepo)

	api := app.Group("/api")
	auth := api.Group("/auth")

	auth.Post("/signup", handlers.SignUpHandler(authService, queries))
	auth.Post("/login", handlers.LoginHandler(authService, queries))

	api.Use(pasetoware.New(pasetoware.Config{
		PrivateKey:  utils.GetPrivateKey(),
		PublicKey:   utils.GetPrivateKey().Public(),
		TokenLookup: [2]string{pasetoware.LookupCookie, "chat_app"},
	}))
}
