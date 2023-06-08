package services

import (
	"chat_backend/generated"
	"chat_backend/internal/app/repositories"
)

type AuthService interface {
	GetUserByUsername(username string) (generated.User, error)
	CreateNewUser(input *repositories.AuthInput) error
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

func (a *authService) CreateNewUser(input *repositories.AuthInput) error {
	return a.authRepository.CreateNewUser(input)
}

func (a *authService) GetUserByUsername(username string) (generated.User, error) {
	return a.authRepository.GetUserByUsername(username)
}

func NewAuthService(r repositories.AuthRepository) AuthService {
	return &authService{
		authRepository: r,
	}
}
