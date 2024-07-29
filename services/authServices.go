package services

import (
	"context"
	"errors"

	"github.com/devmizumizurice/go-jwt-graphql/models"
	"github.com/devmizumizurice/go-jwt-graphql/models/dtos"
	"github.com/devmizumizurice/go-jwt-graphql/repositories"
	"github.com/devmizumizurice/go-jwt-graphql/utils"
)

type AuthServiceInterface interface {
	CreateUser(user *models.User) (*models.User, error)
	VerifyUserEmail(email string, password string) (*dtos.Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dtos.Token, error)
}

type authService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewAuthService(userRepository repositories.UserRepositoryInterface) AuthServiceInterface {
	return &authService{userRepository: userRepository}
}

func (s *authService) CreateUser(user *models.User) (*models.User, error) {
	existingUser, _ := s.userRepository.FindByEmail(user.Email)
	if existingUser != nil {
		return nil, errors.New("EMAIL_ALREADY_EXISTS")
	}

	hash, err := utils.PasswordEncrypt(user.Password)

	if err != nil {
		return nil, errors.New("ERROR_ENCRYPTING_PASSWORD")
	}

	user.Password = hash

	user, err = s.userRepository.Create(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func issueToken(userId string) (*dtos.Token, error) {
	accessToken, err := utils.GenerateToken(userId, false)
	if err != nil {
		return nil, errors.New("ERROR_WHILE_SIGNATURE")
	}
	refreshToken, err := utils.GenerateToken(userId, true)
	if err != nil {
		return nil, errors.New("ERROR_WHILE_SIGNATURE")
	}

	return &dtos.Token{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}

func (s *authService) VerifyUserEmail(email string, password string) (*dtos.Token, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, errors.New("INVALID_EMAIL_OR_PASSWORD")
	}

	err = utils.CompareHashAndPassword(user.Password, password)

	if err != nil {
		return nil, errors.New("INVALID_EMAIL_OR_PASSWORD")
	}

	return issueToken(user.ID)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dtos.Token, error) {
	sub := ctx.Value(utils.SubKey).(string)
	user, err := s.userRepository.FindByID(sub)
	if err != nil {
		return nil, errors.New("USER_NOT_FOUND")
	}

	return issueToken(user.ID)
}
