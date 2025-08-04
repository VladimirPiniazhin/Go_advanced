package auth

import (
	"errors"
	"go/order-api/internals/user"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}
	user := &user.User{
		Email:    email,
		Password: "",
		Name:     name,
	}
	_, err := service.UserRepository.CreateUser(user)
	if err != nil {
		return "", nil
	}
	return user.Email, nil

}
