package services

import (
	"chat_backend/generated"
	"chat_backend/internal/app/repositories"
)

type AuthService interface {
	GetUserByUsername(q *generated.Queries, username string) (generated.User, error)
	CreateNewUser(q *generated.Queries, input *repositories.AuthInput) error
	HashPassword(password string) ([]byte, error)
	VerifyPassword(currentPassword, password string) (bool, error)
}

type authService struct {
	authRepository repositories.AuthRepository
}

func (a *authService) HashPassword(password string) ([]byte, error) {
	return a.authRepository.HashPassword(password)
}

func (a *authService) VerifyPassword(currentPassword, password string) (bool, error) {
	return a.authRepository.VerifyPassword(currentPassword, password)
}

func (a *authService) CreateNewUser(q *generated.Queries, input *repositories.AuthInput) error {
	return a.authRepository.CreateNewUser(q, input)
}

func (a *authService) GetUserByUsername(q *generated.Queries, username string) (generated.User, error) {
	return a.authRepository.GetUserByUsername(q, username)
}

func NewService(r repositories.AuthRepository) AuthService {
	return &authService{
		authRepository: r,
	}
}
