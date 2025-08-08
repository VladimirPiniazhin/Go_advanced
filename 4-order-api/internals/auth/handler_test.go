package auth_test

import (
	"bytes"
	"encoding/json"
	"go/order-api/internals/auth"
	"go/order-api/internals/user"
	"go/order-api/pkg/db"
	"go/order-api/pkg/jwt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	jwtService := jwt.NewJWT(os.Getenv("SECRET"))
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})
	handler := auth.AuthHandler{
		AuthService: auth.NewAuthService(userRepo, jwtService),
	}
	return &handler, mock, nil
}
func TestLoginSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	if err != nil {
		t.Fatal(err)
		return
	}
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("lk@mail.ru", "$2a$10$8cnP5MJNu/OUyTEZhZlxaOl9qE2ohfRHWMJYcho8OKuTBp7ZsqeBa")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "lk@mail.ru",
		Password: "111111",
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("got %d, expected %d", w.Code, 200)
	}

}
func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password", "name", "phone"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	if err != nil {
		t.Fatal(err)
		return
	}

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "lk@mail.ru",
		Password: "111111",
		Name:     "Test_user",
		Phone:    "+79535250880",
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	handler.Register()(w, req)
	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("got %d, expected %d", w.Code, 201)
	}

}
