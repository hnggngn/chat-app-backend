package services

import (
	"chat_backend/generated"
	"chat_backend/internal/app/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserByID(id uuid.UUID) (generated.GetUserByIDRow, error)
	UpdateUser(input *repositories.UpdateInput, id uuid.UUID) error
	DeleteUser(id uuid.UUID) error
}

type userService struct {
	userRepository repositories.UserRepository
}

func (u *userService) DeleteUser(id uuid.UUID) error {
	return u.userRepository.DeleteUser(id)
}

func (u *userService) UpdateUser(input *repositories.UpdateInput, id uuid.UUID) error {
	return u.userRepository.UpdateUser(input, id)
}

func (u *userService) GetUserByID(id uuid.UUID) (generated.GetUserByIDRow, error) {
	return u.userRepository.GetUserByID(id)
}

func NewUserService(r repositories.UserRepository) UserService {
	return &userService{
		userRepository: r,
	}
}
