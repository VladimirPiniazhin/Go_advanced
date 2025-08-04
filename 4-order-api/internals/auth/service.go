package auth

import (
	"errors"
	"go/order-api/internals/link"
	"go/order-api/internals/user"
	"go/order-api/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *user.UserRepository
	jwt            *jwt.JWT
}

func NewAuthService(userRepository *user.UserRepository, j *jwt.JWT) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
		jwt:            j,
	}
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	_, err = service.UserRepository.CreateUser(user)
	if err != nil {
		return "", nil
	}

	token, err := service.jwt.Create(user.Email)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (service *AuthService) UserLogin(email, password string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser == nil {
		return "", errors.New(ErrWrongCredentials)
	}

	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	token, err := service.jwt.Create(email)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (service *AuthService) VerifyByPhone(phoneNumber string) (string, error) {
	existedUser, _ := service.UserRepository.FindByPhoneNumber(phoneNumber)
	if existedUser == nil {
		return "", errors.New(ErrWrongCredentials)
	}
	sessionID := link.RandsStringRunes(10)

	return sessionID, nil
}
