package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `json:"email" validate:"required,email" gorm:"index"`
	Password  string
	Name      string
	Hash      string `json:"hash"`
	Phone     string
	SessionID string `json:"sessionID"`
}

type UsersData struct {
	Users []User `json:"users"`
}
