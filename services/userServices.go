package services

import (
	"errors"

	"github.com/devmizumizurice/go-jwt-graphql/models"
	"github.com/devmizumizurice/go-jwt-graphql/repositories"
)

type UserServiceInterface interface {
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserService(userRepository repositories.UserRepositoryInterface) UserServiceInterface {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userRepository.FindByID(id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, errors.New("USER_NOT_FOUND")
	}

	return user, nil
}
