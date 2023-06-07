package handlers

import (
	"chat_backend/generated"
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

func SignUpHandler(s services.AuthService, q *generated.Queries) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		input := new(repositories.AuthInput)

		if err := ctx.BodyParser(input); err != nil {
			return err
		}

		v := validate.New(input)
		if !v.Validate() {
			return ctx.Status(fiber.StatusForbidden).JSON(v.Errors)
		}

		user, err := s.GetUserByUsername(q, input.Username)
		if err != nil {
			log.Printf("Error in /signup - check user exist: %v", err)
		}

		if user.ID.Valid {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "User already exists.",
			})
		}

		err = s.CreateNewUser(q, input)
		if err != nil {
			log.Printf("Error in /signup - create new user: %v", err)
		}

		return ctx.SendStatus(fiber.StatusCreated)
	}
}

func LoginHandler(s services.AuthService, q *generated.Queries) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		input := new(repositories.AuthInput)

		if err := ctx.BodyParser(input); err != nil {
			return err
		}

		v := validate.New(input)
		if !v.Validate() {
			return ctx.Status(fiber.StatusForbidden).JSON(v.Errors)
		}

		user, err := s.GetUserByUsername(q, input.Username)
		if err != nil {
			log.Printf("Error in /login - check user exist: %v", err)
		}

		if !user.ID.Valid {
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

		userID, _ := user.ID.MarshalJSON()
		token, err := pasetoware.CreateToken(utils.GetPrivateKey(), string(userID), 1*time.Hour, pasetoware.PurposePublic)
		if err != nil {
			log.Printf("Error in /login - create new paseto token: %v", err)
		}

		prod, _ := strconv.ParseBool(os.Getenv("PROD"))

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
