package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Hash     string `json:"hash"`
}

type UsersData struct {
	Users []User `json:"users"`
}
