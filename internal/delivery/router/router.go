package router

import (
	"chat_backend/generated"
	"chat_backend/internal/app/repositories"
	"chat_backend/internal/app/services"
	"chat_backend/internal/delivery/handlers"
	"chat_backend/pkg/utils"
	"github.com/cloudinary/cloudinary-go/v2"
	pasetoware "github.com/gofiber/contrib/paseto"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"
	"time"
)

func AppRouter(app *fiber.App, queries *generated.Queries, cld *cloudinary.Cloudinary) {
	authRepo := repositories.NewAuthRepo(queries)
	authService := services.NewAuthService(authRepo)
	userRepo := repositories.NewUserRepo(queries, cld, authRepo)
	userService := services.NewUserService(userRepo)

	app.Get("/metrics", monitor.New(monitor.Config{
		Title:   "ChatApp Resource Monitor",
		Refresh: 5 * time.Second,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	auth := api.Group("/auth")
	user := api.Group("/user")

	auth.Post("/signup", handlers.SignUpHandler(authService))
	auth.Post("/login", handlers.LoginHandler(authService))

	api.Use(pasetoware.New(pasetoware.Config{
		PrivateKey:  utils.GetPrivateKey(),
		PublicKey:   utils.GetPrivateKey().Public(),
		TokenLookup: [2]string{pasetoware.LookupCookie, "chat_app"},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return fiber.ErrUnauthorized
		},
	}))

	auth.Post("/signout", handlers.SignOutHandler())

	user.Get("/profile", handlers.GetProfileHandler(userService))
	user.Patch("/profile/update", handlers.UpdateProfileHandler(userService))
	user.Delete("/profile/delete", handlers.DeleteUserHandler(userService))
}
