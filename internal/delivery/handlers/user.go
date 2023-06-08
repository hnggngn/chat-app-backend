package handlers

import (
	"chat_backend/internal/app/repositories"
	"chat_backend/internal/app/services"
	"chat_backend/pkg/utils"
	pasetoware "github.com/gofiber/contrib/paseto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gookit/validate"
	"log"
	"mime/multipart"
	"time"
)

func GetProfileHandler(s services.UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID, err := uuid.Parse(ctx.Locals(pasetoware.DefaultContextKey).(string))
		if err != nil {
			log.Printf("Error in /profile - parse uuid: %v", err)
		}

		userData, err := s.GetUserByID(userID)
		if err != nil {
			log.Printf("Error in /profile - get user by id: %v", err)
		}

		return ctx.Status(fiber.StatusOK).JSON(userData)
	}
}

func UpdateProfileHandler(s services.UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		input := new(repositories.UpdateInput)

		if err := ctx.BodyParser(input); err != nil {
			return err
		}

		v := validate.New(input)
		if !v.Validate() {
			return ctx.Status(fiber.StatusForbidden).JSON(v.Errors)
		}

		file, _ := ctx.FormFile("avatar")

		if file != nil {
			if !utils.IsImageFile(file) {
				return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
					"message": "Only image file are allowed (jpeg/png).",
				})
			}

			buffer, _ := file.Open()
			defer func(buffer multipart.File) {
				_ = buffer.Close()
			}(buffer)

			input.Avatar = buffer
		}

		userID, err := uuid.Parse(ctx.Locals(pasetoware.DefaultContextKey).(string))
		if err != nil {
			log.Printf("Error in /profile/update - parse uuid: %v", err)
		}

		err = s.UpdateUser(input, userID)
		if err != nil {
			log.Printf("Error in /profile/update - update user: %v", err)
		}

		return ctx.SendStatus(fiber.StatusOK)
	}
}

func DeleteUserHandler(s services.UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID, err := uuid.Parse(ctx.Locals(pasetoware.DefaultContextKey).(string))
		if err != nil {
			log.Printf("Error in /profile/update - parse uuid: %v", err)
		}

		err = s.DeleteUser(userID)
		if err != nil {
			log.Printf("Error in /profile/delete - delete user: %v", err)
		}

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
