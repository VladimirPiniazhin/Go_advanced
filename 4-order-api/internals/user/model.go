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
	Orders   []Order `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Order struct {
	gorm.Model
	UserID uint `json:"user_id"`
}

type Session struct {
	SessionID string `json:"session_id"`
	Code      string `json:"code"`
}
