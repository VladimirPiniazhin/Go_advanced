package auth

import (
	"errors"
	"go/order-api/internals/link"
	"go/order-api/internals/user"
	"go/order-api/pkg/jwt"
	"math/rand"
	"strconv"

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

func (service *AuthService) Register(email, password, name, phone string) (string, error) {
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
		Phone:    phone,
	}
	_, err = service.UserRepository.CreateUser(user)
	if err != nil {
		return "", nil
	}

	token, err := service.jwt.Create(jwt.JWTData{Phone: phone})
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
	token, err := service.jwt.Create(jwt.JWTData{Phone: email})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (service *AuthService) UpdateSessionID(phoneNumber string) (string, string, error) {
	existedUser, _ := service.UserRepository.FindByPhoneNumber(phoneNumber)
	if existedUser == nil {
		return "", "", errors.New(ErrWrongCredentials)
	}
	existedUser.Session.SessionID = link.RandsStringRunes(10)
	code := strconv.Itoa(rand.Intn(1000))
	existedUser.Session.Code = code
	_, err := service.UserRepository.PatchUser(existedUser)
	if err != nil {
		return "", "", errors.New(ErrInternalError)
	}

	return existedUser.Session.SessionID, code, nil
}

func (service *AuthService) VerifyUserBySmsCode(sessionID, code string) (string, error) {
	existedUser, _ := service.UserRepository.FindBySession(sessionID)
	if existedUser == nil {
		return "", errors.New(ErrWrongCredentials)
	}
	if existedUser.Session.Code != code {
		return "", errors.New(ErrWrongCredentials)
	}

	token, err := service.jwt.Create(jwt.JWTData{Phone: existedUser.Phone})
	if err != nil {
		return "", err
	}
	return token, nil
}
