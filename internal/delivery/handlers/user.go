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

type GetUserByIDRowSwagger struct {
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetProfileHandler retrieves the user profile.
//
//	@Summary		Get user profile
//	@Description	Retrieves the user profile
//	@Tags			Profile
//	@Produce		plain
//	@Success		200	{object}	GetUserByIDRowSwagger
//	@Failure		500	{object}	ErrorResponseSwagger
//	@Router			/user/profile [get]
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

// UpdateProfileHandler updates the user profile.
//
//	@Summary		Update user profile
//	@Description	Updates the user profile
//	@Tags			Profile
//	@Accept			json
//	@Accept			mpfd
//	@Produce		plain
//	@Param			avatar	formData	file					false	"Avatar file (jpeg/png)"
//	@Param			input	body		repositories.AuthInput	false	"User update details"
//	@Success		200		{string}	string					"OK"
//	@Failure		400
//	@Failure		422
//	@Router			/user/profile/update [patch]
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

// DeleteUserHandler deletes the user account.
//
//	@Summary		Delete user account
//	@Description	Deletes the user account
//	@Tags			Profile
//	@Produce		plain
//	@Success		200	{string}	string	"OK"
//	@Failure		500	{object}	ErrorResponseSwagger
//	@Router			/user/profile/delete [delete]
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
