package main

import (
	"bytes"
	"encoding/json"
	"go/order-api/internals/auth"
	"go/order-api/internals/user"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db

}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "lk@mail.ru",
		Password: "$2a$10$8cnP5MJNu/OUyTEZhZlxaOl9qE2ohfRHWMJYcho8OKuTBp7ZsqeBa",
		Name:     "Test_user_1",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", "lk@mail.ru").
		Delete(&user.User{})
}
func TestLoginSuccess(t *testing.T) {
	db := initDb()
	initData(db)
	testServer := httptest.NewServer(App())
	defer testServer.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "lk@mail.ru",
		Password: "111111",
	})

	response, err := http.Post(testServer.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		t.Fatal(err)
	}
	var resData auth.AuthorizationResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	if resData.Token == "" {
		t.Fatal("Token is empty")
	}
	removeData(db)

}

func TestLoginFailed(t *testing.T) {
	testServer := httptest.NewServer(App())
	defer testServer.Close()
	db := initDb()
	initData(db)

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "iva@gmail.com",
		Password: "111111",
	})

	res, err := http.Post(testServer.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal()
	}

	if res.StatusCode == 200 {
		t.Fatalf("Expected error 4XX,  got %d", res.StatusCode)
	}
	removeData(db)
}
