package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" validate:"required,email" gorm:"index"`
	Password string
	Name     string
	Hash     string  `json:"hash"`
	Phone    string  `json:"phone"`
	Session  Session `gorm:"embedded"`
}

type Session struct {
	SessionID string `json:"session_id"`
	Code      string `json:"code"`
}

type UsersData struct {
	Users []User `json:"users"`
}
