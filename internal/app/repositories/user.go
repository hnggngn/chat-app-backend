package repositories

import (
	"chat_backend/generated"
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"mime/multipart"
)

type UserRepository interface {
	GetUserByID(id uuid.UUID) (generated.GetUserByIDRow, error)
	UpdateUser(input *UpdateInput, id uuid.UUID) error
	DeleteUser(id uuid.UUID) error
}

type UpdateInput struct {
	Username string         `form:"username,omitempty" validate:"max_len:30"`
	Password string         `form:"password,omitempty" validate:"max_len:100"`
	Avatar   multipart.File `form:"avatar,omitempty"`
}

type userRepository struct {
	Queries        *generated.Queries
	Cloudinary     *cloudinary.Cloudinary
	AuthRepository AuthRepository
}

func (u *userRepository) DeleteUser(id uuid.UUID) error {
	return u.Queries.DeleteUser(context.Background(), id)
}

func (u *userRepository) UpdateUser(input *UpdateInput, id uuid.UUID) error {
	updates := make(map[string]interface{})

	if input.Avatar != nil {
		upload, err := u.Cloudinary.Upload.Upload(context.Background(), input.Avatar, uploader.UploadParams{
			Folder:         fmt.Sprintf("chat_app/%v/avatar", id),
			PublicID:       "avatar",
			Transformation: "c_crop,g_auto,h_1300,w_1300/f_auto/q_auto:good",
			AllowedFormats: api.CldAPIArray{"png", "jpeg"},
		})
		if err != nil {
			return err
		}

		updates["avatar"] = pgtype.Text{
			String: upload.SecureURL,
			Valid:  true,
		}
	}

	if len(input.Password) > 0 {
		password, err := u.AuthRepository.HashPassword(input.Password)
		if err != nil {
			return err
		}

		updates["password"] = string(password)
	}

	if len(input.Username) > 0 {
		updates["username"] = input.Username
	}

	return u.Queries.UpdateUser(context.Background(), generated.UpdateUserParams{
		Column1: updates["username"],
		Column2: updates["password"],
		Column3: updates["avatar"],
		ID:      id,
	})
}

func (u *userRepository) GetUserByID(id uuid.UUID) (generated.GetUserByIDRow, error) {
	return u.Queries.GetUserByID(context.Background(), id)
}

func NewUserRepo(queries *generated.Queries, cld *cloudinary.Cloudinary, repository AuthRepository) UserRepository {
	return &userRepository{
		Queries:        queries,
		Cloudinary:     cld,
		AuthRepository: repository,
	}
}
