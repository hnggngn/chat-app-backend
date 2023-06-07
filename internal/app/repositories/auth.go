package repositories

import (
	"chat_backend/generated"
	"chat_backend/pkg/utils"
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/matthewhartstonge/argon2"
	"log"
)

type AuthRepository interface {
	GetUserByUsername(q *generated.Queries, username string) (generated.User, error)
	CreateNewUser(q *generated.Queries, input *AuthInput) error
	HashPassword(password string) ([]byte, error)
	VerifyPassword(currentPassword, password string) (bool, error)
}

type authRepository struct {
	Queries *generated.Queries
}

func (r *authRepository) HashPassword(password string) ([]byte, error) {
	argon := argon2.DefaultConfig()
	return argon.HashEncoded([]byte(password))
}

func (r *authRepository) VerifyPassword(currentPassword, password string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(currentPassword))
}

type AuthInput struct {
	Username string `json:"username" validate:"required|max_len:30"`
	Password string `json:"password" validate:"required|max_len:100"`
}

func (r *authRepository) CreateNewUser(q *generated.Queries, input *AuthInput) error {
	hash, err := r.HashPassword(input.Password)
	if err != nil {
		return err
	}

	err = q.CreateNewUser(context.Background(), generated.CreateNewUserParams{
		Username: input.Username,
		Password: string(hash),
		Avatar: pgtype.Text{
			String: utils.GetRandomImage(),
			Valid:  true,
		},
	})
	if err != nil {
		return err
	}
	if err != nil {
		log.Printf("Error in /signup - create new user: %v", err)
	}

	return nil
}

func (r *authRepository) GetUserByUsername(q *generated.Queries, username string) (generated.User, error) {
	return q.GetUserByUsername(context.Background(), username)
}

func NewRepo(queries *generated.Queries) AuthRepository {
	return &authRepository{
		Queries: queries,
	}
}
