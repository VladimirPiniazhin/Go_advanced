package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" validate:"required,email" gorm:"index"`
	Password string
	Name     string
	Hash     string `json:"hash"`
}

type UsersData struct {
	Users []User `json:"users"`
}
