package di

import "go/order-api/internals/user"

type IStatRepository interface {
	AddClick(linkId uint)
}

type IUserRepository interface {
	CreateUser(user *user.User) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
	FindByPhoneNumber(phone string) (*user.User, error)
	FindBySession(session string) (*user.User, error)
	PatchUser(user *user.User) (*user.User, error)
	DeleteUser(id uint) error
}
