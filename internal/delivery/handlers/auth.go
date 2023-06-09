package handlers

import (
	"chat_backend/internal/app/repositories"
	"chat_backend/internal/app/services"
	"chat_backend/pkg/utils"
	pasetoware "github.com/gofiber/contrib/paseto"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	prod, _ = strconv.ParseBool(os.Getenv("PROD"))
)

type ErrorResponseSwagger struct {
	Message string `json:"message"`
}

// SignUpHandler handles the signup route.
//
//	@Summary		Handle user registration and create a new user.
//	@Description	Handle user registration and create a new user.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		plain
//	@Param			input	body		repositories.AuthInput	true	"User registration details"
//	@Success		201		{string}	string					"Created"
//	@Failure		400		{object}	ErrorResponseSwagger
//	@Failure		403
//	@Router			/auth/signup [post]
func SignUpHandler(s services.AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		input := new(repositories.AuthInput)

		if err := ctx.BodyParser(input); err != nil {
			return err
		}

		v := validate.New(input)
		if !v.Validate() {
			return ctx.Status(fiber.StatusForbidden).JSON(v.Errors)
		}

		user, _ := s.GetUserByUsername(input.Username)

		if len(user.Username) > 0 {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "User already exists.",
			})
		}

		err := s.CreateNewUser(input)
		if err != nil {
			log.Printf("Error in /signup - create new user: %v", err)
		}

		return ctx.SendStatus(fiber.StatusCreated)
	}
}

// LoginHandler handles the login route.
//
//	@Summary		Handle user login and generate an authentication token.
//	@Description	Handle user login and generate an authentication token.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		plain
//	@Param			input	body		repositories.AuthInput	true	"User login details"
//	@Success		200		{string}	string					"OK"
//	@Failure		400		{object}	ErrorResponseSwagger
//	@Failure		403
//	@Router			/auth/login [post]
func LoginHandler(s services.AuthService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		input := new(repositories.AuthInput)

		if err := ctx.BodyParser(input); err != nil {
			return err
		}

		v := validate.New(input)
		if !v.Validate() {
			return ctx.Status(fiber.StatusForbidden).JSON(v.Errors)
		}

		user, _ := s.GetUserByUsername(input.Username)

		if len(user.Username) == 0 {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "User not exists.",
			})
		}

		verifyPassword, err := s.VerifyPassword(user.Password, input.Password)
		if err != nil {
			log.Printf("Error in /login - verify password: %v", err)
		}
		if !verifyPassword {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Password not correct.",
			})
		}

		userID := user.ID.String()
		token, err := pasetoware.CreateToken(utils.GetPrivateKey(), userID, 1*time.Hour, pasetoware.PurposePublic)
		if err != nil {
			log.Printf("Error in /login - create new paseto token: %v", err)
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     "chat_app",
			Value:    token,
			HTTPOnly: prod,
			Secure:   prod,
			SameSite: fiber.CookieSameSiteStrictMode,
		})

		return ctx.SendStatus(fiber.StatusOK)
	}
}

// SignOutHandler handles the signout route.
//
//	@Summary		Handle user signout and remove the authentication token.
//	@Description	Handle user signout and remove the authentication token.
//	@Tags			Authentication
//	@Produce		plain
//	@Success		200	{string}	string	"OK"
//	@Router			/auth/signout [post]
func SignOutHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Cookie(&fiber.Cookie{
			Name:     "chat_app",
			Value:    "",
			HTTPOnly: prod,
			Secure:   prod,
			SameSite: fiber.CookieSameSiteStrictMode,
			Expires:  time.Now().Add(-time.Hour),
		})
		return ctx.SendStatus(fiber.StatusOK)
	}
}
