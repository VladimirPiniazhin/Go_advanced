package jwt_test

import (
	"go/order-api/pkg/jwt"
	"os"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	const email = "lk@mail.ru"

	jwtService := jwt.NewJWT(os.Getenv("SECRET"))
	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatalf("Invalid token")
	}
	if data.Email != email {
		t.Fatalf("Email %s is not equal to %s", data.Email, email)
	}
}
