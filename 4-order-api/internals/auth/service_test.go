package auth_test

import (
	"go/order-api/internals/auth"
	"go/order-api/internals/user"
	"go/order-api/pkg/jwt"
	"os"
	"testing"
)

type MockUserRepository struct {
}

func (repo *MockUserRepository) CreateUser(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "lk@mail.ru",
	}, nil
}
func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}
func (repo *MockUserRepository) FindByPhoneNumber(phone string) (*user.User, error) {
	return nil, nil
}
func (repo *MockUserRepository) FindBySession(session string) (*user.User, error) {
	return nil, nil
}
func (repo *MockUserRepository) PatchUser(user *user.User) (*user.User, error) {
	return nil, nil
}
func (repo *MockUserRepository) DeleteUser(id uint) error {
	return nil
}

func TestRegisterSuccess(t *testing.T) {
	const (
		Email    = "lk@mail.ru"
		Password = "111111"
		Name     = "Test_user_1"
		Phone    = "+79535250880"
		Token    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImxrQG1haWwucnUifQ.z-LbhTtgfgnN0tniUZdQ-l5oDPbHSmHs5MJknC3uKBc"
	)
	jwtService := jwt.NewJWT(os.Getenv("SECRET"))
	authService := auth.NewAuthService(&MockUserRepository{}, jwtService)
	token, err := authService.Register(Email, Password, Name, Phone)
	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is not valid")
	}

	if data.Email != Email {
		t.Fatalf("Email %s is not equal to %s", data.Email, Email)
	}

}
